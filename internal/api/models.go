package api

// SubjectSummary is a summary of a subject for list responses.
type SubjectSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SubjectCreateRequest is the request body for POST /v1/subjects.
type SubjectCreateRequest struct {
	SubjectName string `json:"subject_name"`
}

// SubjectUpdateRequest is the request body for PUT /v1/subjects/{subjectID}.
type SubjectUpdateRequest struct {
	SubjectName string `json:"subject_name"`
}

// StatusResponse is a generic response for mutation operations (create, update, delete).
type StatusResponse struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}
