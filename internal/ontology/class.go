package ontology

const ClassDType = "class"

type Class struct {
	ID            string         `json:"uid"`
	Name          string         `json:"name"`
	IsRoot        bool           `json:"isRoot"`
	PropertyTypes []PropertyType `json:"propertyTypes,omitempty"`
	Individuals   []Individual   `json:"individuals,omitempty"`
	DType         []string       `json:"dgraph.type,omitempty"`
	Subclasses    []Class        `json:"subclasses,omitempty"`
	IsSoftDeleted bool           `json:"isSoftDeleted"`
}

type ClassUpdate struct {
	ID            string         `json:"uid"`
	Name          *string        `json:"name,omitempty"`
	IsRoot        *bool          `json:"isRoot,omitempty"`
	PropertyTypes []PropertyType `json:"propertyTypes,omitempty"`
	Individuals   []Individual   `json:"individuals,omitempty"`
	DType         []string       `json:"dgraph.type,omitempty"`
	Subclasses    []Class        `json:"subclasses,omitempty"`
	IsSoftDeleted *bool          `json:"isSoftDeleted,omitempty"`
}
