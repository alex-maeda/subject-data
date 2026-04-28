package datamodel

// Subject represents a subject.
type Subject struct {
	Identifiable
	SubjectName string `json:"subject_name" jsonschema_description:"The subject's display name (e.g. 'Last, First')"`
}
