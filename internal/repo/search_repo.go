package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v200"
)

type Search struct {
	client *dgo.Dgraph
}

func NewSearch(client *dgo.Dgraph) *Search {
	return &Search{client: client}
}

type SearchFilter struct {
	Name          *string
	PropertyValue *string
	PropertyType  *string
	PropertyName  *string
}

type SearchResponse struct {
	ID   string   `json:"uid"`
	Name string   `json:"name"`
	Type []string `json:"dgraph.type"`
}

func (s *Search) Query(ctx context.Context, filter SearchFilter) ([]SearchResponse, error) {
	rootFunc := "has(dgraph.type)"
	if filter.Name != nil {
		rootFunc = fmt.Sprintf("regexp(name,/.*%s.*/)", *filter.Name)
	}
	filterPropertyValue := ""
	if filter.PropertyValue != nil {
		filterPropertyValue = fmt.Sprintf("@filter(eq(value,\"%s\"))", *filter.PropertyValue)
	}
	filterPropertyType := ""
	if filter.PropertyType != nil && filter.PropertyName != nil {
		filterPropertyType = fmt.Sprintf("@filter(eq(propertyType,\"%s\") and eq(name,\"%s\"))", *filter.PropertyType, *filter.PropertyName)
	} else if filter.PropertyType != nil {
		filterPropertyType = fmt.Sprintf("@filter(eq(propertyType,\"%s\"))", *filter.PropertyType)
	} else if filter.PropertyName != nil {
		filterPropertyType = fmt.Sprintf("@filter(eq(name,\"%s\"))", *filter.PropertyName)
	}
	cascade := "@cascade"
	if filter.PropertyName == nil && filter.PropertyType == nil && filter.PropertyValue == nil {
		cascade = ""
	}
	cascade = "@filter(eq(dgraph.type,class) or eq(dgraph.type,individual))" + cascade
	q := fmt.Sprintf(`query Search(){
  search(func: %s)%s{
    	uid
		name
      propertyValues%s{
		uid
		value
      propertyTypeRef%s{
		uid
		propertyType
        name
      }
    }
    dgraph.type
  }
}`, rootFunc, cascade, filterPropertyValue, filterPropertyType)
	query, err := s.client.NewTxn().Query(ctx, q)
	if err != nil {
		return nil, err
	}
	type resp struct {
		Res []SearchResponse `json:"search"`
	}
	var res resp
	err = json.Unmarshal(query.Json, &res)
	if err != nil {
		return nil, err
	}
	return res.Res, nil
}
