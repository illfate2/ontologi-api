package repo

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"

	"github.com/illfate2/ontology-api/internal/ontology"
)

type Property struct {
	client *dgo.Dgraph
}

func NewProperty(client *dgo.Dgraph) *Property {
	return &Property{client: client}
}

func (p *Property) AddPropertyTypeToClass(ctx context.Context, classUID string, property ontology.PropertyType) (string, error) {
	var parent ontology.ClassUpdate
	parent.ID = classUID
	parent.DType = []string{ontology.ClassDType}
	property.ID = "_:" + property.Name
	property.DType = []string{ontology.PropertyTypeDType}
	parent.PropertyTypes = []ontology.PropertyType{property}
	classJSON, err := json.Marshal(parent)
	if err != nil {
		return "", err
	}
	res, err := p.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   classJSON,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}
	return res.GetUids()[property.Name], err
}

func (p *Property) UpdatePropertyType(ctx context.Context, property ontology.PropertyType) error {
	propertyJSON, err := json.Marshal(property)
	if err != nil {
		return err
	}
	_, err = p.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   propertyJSON,
		CommitNow: true,
	})
	return err
}

func (p *Property) AddPropertyValueToIndividual(ctx context.Context, propertyValue ontology.PropertyValueOneType, individualID string) (string, error) {
	var individual ontology.Individual
	individual.ID = individualID
	individual.DType = []string{ontology.IndividualDType}
	propertyValue.ID = "_:" + propertyValue.Value
	propertyValue.DType = []string{ontology.PropertyValueDType}
	propertyType := propertyValue.Type
	propertyValue.Type = nil
	propertyJSON, _ := json.Marshal(propertyValue)
	res, err := p.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   propertyJSON,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}

	propertyID := res.GetUids()[propertyValue.Value]
	updateProperty := map[string]interface{}{
		"uid": propertyID,
		"propertyTypeRef": map[string]interface{}{
			"uid": propertyType.ID,
		},
	}
	propertyJSON, _ = json.Marshal(updateProperty)
	res, err = p.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   propertyJSON,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}

	propertyValue.Value = ""
	propertyValue.Type = nil
	propertyValue.DType = nil
	propertyValue.ID = propertyID
	individual.PropertyValues = []ontology.PropertyValue{ontology.PropertyValue{
		ID: propertyID,
	}}
	individualJSON, _ := json.Marshal(individual)
	res, err = p.client.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   individualJSON,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}
	return propertyID, nil

}
