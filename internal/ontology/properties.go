package ontology

const PropertyTypeDType = "property_type"
const PropertyValueDType = "property_value"

type PropertyType struct {
	ID    string   `json:"uid"`
	Name  string   `json:"name"`
	Type  string   `json:"propertyType"`
	DType []string `json:"dgraph.type,omitempty"`
}

type PropertyValue struct {
	ID    string          `json:"uid"`
	Type  []*PropertyType `json:"propertyTypeRef,omitempty"`
	Value string          `json:"value,omitempty"`
	DType []string        `json:"dgraph.type,omitempty"`
}
type PropertyValueOneType struct {
	ID    string        `json:"uid"`
	Type  *PropertyType `json:"propertyType,omitempty"`
	Value string        `json:"value,omitempty"`
	DType []string      `json:"dgraph.type,omitempty"`
}
