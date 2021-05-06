package ontology

const RelationshipDType = "relationship"

type Relationship struct {
	ID    string   `json:"uid"`
	Name  string   `json:"name,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

const RelationshipTripleDType = "relationship_triple"

type RelationshipTriple struct {
	ID      string       `json:"uid"`
	Name    Relationship `json:"relationshipName"`
	Object  Individual   `json:"object"`
	Subject Individual   `json:"subject"`
	DType   []string     `json:"dgraph.type,omitempty"`
}
