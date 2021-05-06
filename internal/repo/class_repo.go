package repo

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"

	"github.com/illfate2/ontology-api/internal/ontology"
)

type Class struct {
	client *dgo.Dgraph
}

func NewClass(client *dgo.Dgraph) *Class {
	return &Class{client: client}
}

func (c *Class) Insert(ctx context.Context, class ontology.Class) (string, error) {
	class.ID = "_:" + class.Name
	class.DType = []string{ontology.ClassDType}

	classJSON, err := json.Marshal(class)
	if err != nil {
		return "", err
	}
	res, err := c.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   classJSON,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}
	return res.GetUids()[class.Name], err
}

func (c *Class) AddSubclassToParent(ctx context.Context, parent, subclass ontology.Class) (string, error) {
	parent.DType = []string{ontology.ClassDType}
	subclass.ID = "_:" + subclass.Name
	subclass.DType = []string{ontology.ClassDType}
	parent.Subclasses = []ontology.Class{subclass}
	classJSON, err := json.Marshal(parent)
	if err != nil {
		return "", err
	}
	res, err := c.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   classJSON,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}
	return res.GetUids()[subclass.Name], err
}

func (c *Class) Find(ctx context.Context, uid string) (ontology.Class, error) {
	q := `query Class($dgraphType: string, $id: string){
  class(func: uid($id),first:1)@filter(eq(dgraph.type,"class")) {
	uid
    name
	isRoot
	isSoftDeleted
	propertyTypes{
		uid
		name
		propertyType
	}
	individuals{
		uid
		name
		isSoftDeleted
		propertyValues{
    		propertyType{
				name
    			isSoftDeleted
    			propertyType
			}
    		isSoftDeleted
    		value
		}
	}
  }
}`
	query, err := c.client.NewTxn().QueryWithVars(ctx, q, map[string]string{
		dgraphType: ontology.ClassDType,
		"$id":      uid,
	})
	if err != nil {
		return ontology.Class{}, err
	}
	type resp struct {
		Classes []ontology.Class `json:"class"`
	}
	var classResp resp
	err = json.Unmarshal(query.Json, &classResp)
	if err != nil {
		return ontology.Class{}, err
	}
	return classResp.Classes[0], nil
}

const (
	isRootField = "$isRoot"
	dgraphType  = "$dgraphType"
)

func (c *Class) FindAll(ctx context.Context) ([]ontology.Class, error) {
	q := `query Classes($dgraphType: string,$isRoot: bool){
  classes(func: eq(dgraph.type,"class"))@recurse@filter(eq(isRoot,$isRoot)) {
	uid
    name
	isRoot
	isSoftDeleted
	subclasses
  }
}`
	query, err := c.client.NewTxn().QueryWithVars(ctx, q, map[string]string{
		isRootField: "true",
		dgraphType:  ontology.ClassDType,
	})
	if err != nil {
		return nil, err
	}
	type resp struct {
		Classes []ontology.Class `json:"classes"`
	}
	var classResp resp
	err = json.Unmarshal(query.Json, &classResp)
	if err != nil {
		return nil, err
	}
	return classResp.Classes, nil
}

func (c *Class) Remove(ctx context.Context, id string) error {
	class := ontology.Class{
		ID:            id,
		IsSoftDeleted: true,
	}
	classJSON, err := json.Marshal(class)
	if err != nil {
		return err
	}
	_, err = c.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   classJSON,
		CommitNow: true,
	})
	return err
}
