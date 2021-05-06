package repo

import (
	"context"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
)

const ontologySchema = `
isRoot: bool .
isSoftDeleted: bool .
Subclasses: [uid] .
propertyType: string .
PropertyTypes: [uid] .

type Class{
    name: String
    isRoot: Bool
    isSoftDeleted: Bool
    Subclasses: [Class]
	Individuals: [Individuals]
	PropertyTypes: [PropertyType]
}

type Relationship{
	name: String 
	isSoftDeleted: Bool
}

relationshipName: uid .
subject: uid .
object: uid .

type RelationshipTriple{
	relationshipName:  Relationship
	subject: Individual
	object: Individual
	isSoftDeleted: Bool
}

type PropertyType {
    name: String
    isSoftDeleted: Bool
    propertyType: String
}


value: String .
propertyTypeRef: [uid] .


type PropertyValue {
    propertyTypeRef: PropertyType
    isSoftDeleted: Bool
    value: String
}

PropertyValues: [uid] .
Individuals: [uid] .
name: String @index(term, fulltext, trigram) .
relationshipTriples: [uid] .

type Individual {
    PropertyValues: [PropertyValue]
    name: String
    isSoftDeleted: Bool
	relationshipTriples: [RelationshipTriple]
}

`

func MigrateSchema(ctx context.Context, dgraph *dgo.Dgraph) error {
	return dgraph.Alter(ctx, &api.Operation{
		Schema: ontologySchema,
	})
}
