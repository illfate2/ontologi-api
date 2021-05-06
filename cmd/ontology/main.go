package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/illfate2/ontology-api/internal/config"
	"github.com/illfate2/ontology-api/internal/repo"
	"github.com/illfate2/ontology-api/internal/server/rest"
)

const (
	dgraphURIEnv = "DGRAPH_URI"
	serverPort   = "SERVER_PORT"
)

func main() {
	ctx, cancelF := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelF()
	conn, closeF := config.MustGetDgraphConn(ctx, os.Getenv(dgraphURIEnv))
	defer closeF()
	classRepo := repo.NewClass(conn)
	propertyRepo := repo.NewProperty(conn)
	individualRepo := repo.NewIndividual(conn)
	relationshipRepo := repo.NewRelationship(conn)
	searchRepo := repo.NewSearch(conn)
	err := repo.MigrateSchema(context.Background(), conn)
	if err != nil {
		panic(err)
	}
	server := rest.NewServer(classRepo, propertyRepo, individualRepo, relationshipRepo, searchRepo)
	log.Panic(http.ListenAndServe(":"+os.Getenv(serverPort), server))
}
