package ontology

const IndividualDType = "individual"

type Individual struct {
	ID             string               `json:"uid"`
	Name           *string              `json:"name,omitempty"`
	PropertyValues []PropertyValue      `json:"propertyValues,omitempty"`
	DType          []string             `json:"dgraph.type,omitempty"`
	Relationships  []RelationshipTriple `json:"relationshipTriples,omitempty"`
}
