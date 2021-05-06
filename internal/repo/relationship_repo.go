package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"

	"github.com/illfate2/ontology-api/internal/ontology"
)

type Relationship struct {
	client *dgo.Dgraph
}

func NewRelationship(client *dgo.Dgraph) *Relationship {
	return &Relationship{client: client}
}

func (r *Relationship) Create(ctx context.Context, relationship ontology.Relationship) (string, error) {
	relationship.ID = "_:" + relationship.Name
	relationship.DType = []string{ontology.RelationshipDType}

	relationshipJSON, err := json.Marshal(relationship)
	if err != nil {
		return "", err
	}
	res, err := r.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   relationshipJSON,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}
	return res.GetUids()[relationship.Name], err
}

func (r *Relationship) Update(ctx context.Context, relationship ontology.Relationship) error {
	relationship.DType = []string{ontology.RelationshipDType}

	relationshipJSON, err := json.Marshal(relationship)
	if err != nil {
		return err
	}
	_, err = r.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   relationshipJSON,
		CommitNow: true,
	})
	return err
}

func (r *Relationship) Delete(ctx context.Context, id string) error {
	relationship := ontology.Relationship{
		ID:    id,
		DType: []string{ontology.RelationshipDType},
	}

	relationshipJSON, err := json.Marshal(relationship)
	if err != nil {
		return err
	}
	_, err = r.client.NewTxn().Mutate(ctx, &api.Mutation{
		DeleteJson: relationshipJSON,
		CommitNow:  true,
	})
	return err
}

func (r *Relationship) FindAll(ctx context.Context, name *string) ([]ontology.Relationship, error) {
	filterQuery := "/.*.*/"
	if name != nil {
		filterQuery = "/.*" + *name + ".*/"
	}
	q := `query Relationships($dgraphType: string){
  relationships(func: eq(dgraph.type,$dgraphType))@filter(regexp(name,%s)){
	uid
  name
  }
}`
	query, err := r.client.NewTxn().QueryWithVars(ctx, fmt.Sprintf(q, filterQuery), map[string]string{
		dgraphType: ontology.RelationshipDType,
	})
	if err != nil {
		return nil, err
	}
	type resp struct {
		Relationships []ontology.Relationship `json:"relationships"`
	}
	var relResp resp
	err = json.Unmarshal(query.Json, &relResp)
	if err != nil {
		return nil, err
	}
	return relResp.Relationships, nil
}

func (r *Relationship) AddTripleToIndividual(ctx context.Context, individualID string, triple ontology.RelationshipTriple) (string, error) {
	id, err := r.createTriple(ctx, triple)
	if err != nil {
		return "", err
	}
	err = r.addTripleToIndividual(ctx, individualID, id)
	return id, err
}

func (r *Relationship) createTriple(ctx context.Context, triple ontology.RelationshipTriple) (string, error) {
	triple.ID = "_:" + triple.Name.ID
	triple.DType = []string{ontology.RelationshipTripleDType}

	tripleJSON, err := json.Marshal(triple)
	if err != nil {
		return "", err
	}
	res, err := r.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   tripleJSON,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}
	return res.GetUids()[triple.Name.ID], err
}

func (r *Relationship) addTripleToIndividual(ctx context.Context, individualID, tripleID string) error {
	updateReq := map[string]interface{}{
		"uid": individualID,
		"relationshipTriples": map[string]string{
			"uid": tripleID,
		},
	}

	tripleJSON, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}
	_, err = r.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   tripleJSON,
		CommitNow: true,
	})
	return err
}
