package rest

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/illfate2/ontology-api/internal/ontology"
	"github.com/illfate2/ontology-api/internal/repo"
)

type Server struct {
	http.Handler
	classRepo        *repo.Class
	propertyRepo     *repo.Property
	individualRepo   *repo.Individual
	relationshipRepo *repo.Relationship
	searchRepo       *repo.Search
	engine           *gin.Engine
}

func NewServer(classRepo *repo.Class, propertyRepo *repo.Property, individualRepo *repo.Individual, relationshipRepo *repo.Relationship, searchRepo *repo.Search) *Server {
	engine := gin.New()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PATCH", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	s := &Server{
		Handler:          engine,
		engine:           engine,
		classRepo:        classRepo,
		individualRepo:   individualRepo,
		relationshipRepo: relationshipRepo,
		propertyRepo:     propertyRepo,
		searchRepo:       searchRepo,
	}
	engine.GET("/classes", s.HandleGetClasses)
	engine.GET("/search", s.HandleSearch)
	engine.GET("/class/:id", s.HandleGetClass)
	engine.POST("/class", s.HandleCreateClass)
	engine.DELETE("/class/:id", s.HandleDeleteClass)

	engine.POST("/class/:id/property", s.HandleAddPropertyTypeToClass)
	engine.PATCH("/property/:id", s.HandleUpdateProperty)

	engine.POST("/class/:id/individual", s.HandleAddIndividualToClass)
	engine.PATCH("/individual/:id", s.HandleUpdateIndividual)
	engine.GET("/individuals", s.HandleGetIndividuals)

	engine.POST("/relationship", s.HandleCreateRelationship)
	engine.GET("/relationships", s.HandleGetRelationships)
	engine.PATCH("/relationship/:id", s.HandleUpdateRelationship)
	engine.DELETE("/relationship/:id", s.HandleDeleteRelationship)

	engine.POST("/individual/:id/property", s.HandleAddPropertyToIndividual)
	engine.POST("/individual/:id/relationship", s.HandleAddRelationshipTripleToIndividual)

	return s
}

func (s *Server) HandleGetClasses(ctx *gin.Context) {
	all, err := s.classRepo.FindAll(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, all)
}

func (s *Server) HandleGetClass(ctx *gin.Context) {
	all, err := s.classRepo.Find(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, all)
}

func (s *Server) HandleDeleteClass(ctx *gin.Context) {
	id := ctx.Param("id")
	err := s.classRepo.Remove(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (s *Server) HandleDeleteRelationship(ctx *gin.Context) {
	id := ctx.Param("id")
	err := s.relationshipRepo.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}

type createClassReq struct {
	Parent   ontology.Class  `json:"parent"`
	Subclass *ontology.Class `json:"subclass"`
}

func (s *Server) HandleCreateClass(ctx *gin.Context) {
	var req createClassReq
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	var id string
	if req.Subclass == nil {
		id, err = s.classRepo.Insert(ctx.Request.Context(), req.Parent)
	} else {
		id, err = s.classRepo.AddSubclassToParent(ctx.Request.Context(), req.Parent, *req.Subclass)
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"uid": id})
}

func (s *Server) HandleAddPropertyTypeToClass(ctx *gin.Context) {
	var propertyType ontology.PropertyType
	err := ctx.BindJSON(&propertyType)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	propertyID, err := s.propertyRepo.AddPropertyTypeToClass(ctx.Request.Context(), ctx.Param("id"), propertyType)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"uid": propertyID})
}

func (s *Server) HandleUpdateProperty(ctx *gin.Context) {
	var propertyType ontology.PropertyType
	err := ctx.BindJSON(&propertyType)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	propertyType.ID = ctx.Param("id")
	err = s.propertyRepo.UpdatePropertyType(ctx.Request.Context(), propertyType)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func (s *Server) HandleAddIndividualToClass(ctx *gin.Context) {
	var individual ontology.Individual
	err := ctx.BindJSON(&individual)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	individualID, err := s.individualRepo.AddIndividualToClass(ctx.Request.Context(), ctx.Param("id"), individual)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"uid": individualID})
}

func (s *Server) HandleAddPropertyToIndividual(ctx *gin.Context) {
	var value ontology.PropertyValueOneType
	err := ctx.BindJSON(&value)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	propertyID, err := s.propertyRepo.AddPropertyValueToIndividual(ctx.Request.Context(), value, ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"uid": propertyID})
}

func (s *Server) HandleAddRelationshipTripleToIndividual(ctx *gin.Context) {
	var value ontology.RelationshipTriple
	err := ctx.BindJSON(&value)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	tripleID, err := s.relationshipRepo.AddTripleToIndividual(ctx.Request.Context(), ctx.Param("id"), value)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"uid": tripleID})
}

func (s *Server) HandleUpdateIndividual(ctx *gin.Context) {
	var individual ontology.Individual
	err := ctx.BindJSON(&individual)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	individual.ID = ctx.Param("id")
	err = s.individualRepo.UpdateIndividual(ctx.Request.Context(), individual)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func (s *Server) HandleCreateRelationship(ctx *gin.Context) {
	var req ontology.Relationship
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	id, err := s.relationshipRepo.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"uid": id})
}

func (s *Server) HandleUpdateRelationship(ctx *gin.Context) {
	var relationship ontology.Relationship
	err := ctx.BindJSON(&relationship)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	relationship.ID = ctx.Param("id")
	err = s.relationshipRepo.Update(ctx.Request.Context(), relationship)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func (s *Server) HandleGetRelationships(ctx *gin.Context) {
	queryName := ctx.Query("name")
	filter := &queryName
	if queryName == "" {
		filter = nil
	}
	all, err := s.relationshipRepo.FindAll(ctx.Request.Context(), filter)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, all)
}

func (s *Server) HandleGetIndividuals(ctx *gin.Context) {
	queryName := ctx.Query("name")
	queryID := ctx.Query("id")
	filterName := &queryName
	if queryName == "" {
		filterName = nil
	}
	filterID := &queryID
	if queryID == "" {
		filterID = nil
	}
	all, err := s.individualRepo.FindAll(ctx.Request.Context(), filterName, filterID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, all)
}
func (s *Server) HandleSearch(ctx *gin.Context) {
	all, err := s.searchRepo.Query(ctx, repo.SearchFilter{
		Name:          getPtr(ctx.Query("name")),
		PropertyValue: getPtr(ctx.Query("propertyValue")),
		PropertyType:  getPtr(ctx.Query("propertyType")),
		PropertyName:  getPtr(ctx.Query("propertyName")),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, all)
}

func getPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
