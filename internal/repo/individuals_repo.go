package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"

	"github.com/illfate2/ontology-api/internal/ontology"
)

type Individual struct {
	client *dgo.Dgraph
}

func NewIndividual(client *dgo.Dgraph) *Individual {
	return &Individual{client: client}
}

func (i *Individual) FindAll(ctx context.Context, name, id *string) ([]ontology.Individual, error) {
	filterQuery := "/.*.*/"
	if name != nil {
		filterQuery = "/.*" + *name + ".*/"
	}

	q := fmt.Sprintf(`query Individuals($dgraphType: string){
  individuals(func: eq(dgraph.type,$dgraphType))@filter(regexp(name,%s)){
	uid
  	name
  }
}`, filterQuery)
	if id != nil {
		q = fmt.Sprintf(`query Individuals($dgraphType: string){
  individuals(func: uid(%s)){
	uid
  	name
	propertyValues{
		uid
		value
		propertyTypeRef{
			uid
			name
			propertyType
		}
	}
	relationshipTriples{
		uid	
		relationshipName{
			uid
			name
		}
		subject{
			uid
  			name
		}
		object{
			uid
  			name
		}
	}
  }
}`, *id)
	}
	query, err := i.client.NewTxn().QueryWithVars(ctx, q, map[string]string{
		dgraphType: ontology.IndividualDType,
	})
	if err != nil {
		return nil, err
	}
	type resp struct {
		Individuals []ontology.Individual `json:"individuals"`
	}
	var iResp resp
	err = json.Unmarshal(query.Json, &iResp)
	if err != nil {
		return nil, err
	}
	return iResp.Individuals, nil
}

func (i *Individual) AddIndividualToClass(ctx context.Context, classUID string, individual ontology.Individual) (string, error) {
	var parent ontology.ClassUpdate
	parent.ID = classUID
	parent.DType = []string{ontology.ClassDType}

	individual.ID = "_:" + *individual.Name
	individual.DType = []string{ontology.IndividualDType}
	parent.Individuals = []ontology.Individual{individual}
	classJSON, err := json.Marshal(parent)
	if err != nil {
		return "", err
	}
	res, err := i.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   classJSON,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}
	return res.GetUids()[*individual.Name], err
}

func (i *Individual) UpdateIndividual(ctx context.Context, individual ontology.Individual) error {
	individualJSON, err := json.Marshal(individual)
	if err != nil {
		return err
	}
	_, err = i.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   individualJSON,
		CommitNow: true,
	})
	return err
}
