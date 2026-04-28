package datamodel

// Identifiable adds a standard ID field to structs that can be persisted.
// Embed this in structs that need an opaque identifier.
type Identifiable struct {
	ID *string `json:"id,omitempty" jsonschema_description:"Opaque identifier used to reference this object"`
}
