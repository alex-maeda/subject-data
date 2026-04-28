package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ListSubjects handles GET /v1/subjects.
//
//	@Summary	List all subjects
//	@Tags		subjects
//	@Produce	json
//	@Success	200	{array}		SubjectSummary
//	@Failure	500	{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/subjects [get]
func (h *Handler) ListSubjects(w http.ResponseWriter, _ *http.Request) {
	subjects, err := h.SubjectService.ListSubjects()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	items := make([]SubjectSummary, 0, len(subjects))
	for _, s := range subjects {
		item := SubjectSummary{}
		if s.ID != nil {
			item.ID = *s.ID
		}
		if s.SubjectName != nil {
			item.Name = *s.SubjectName
		}
		items = append(items, item)
	}
	h.writeJSON(w, http.StatusOK, items)
}

// GetSubject handles GET /v1/subjects/{subjectID}.
//
//	@Summary	Get subject by ID
//	@Tags		subjects
//	@Produce	json
//	@Param		subjectID	path		string	true	"Subject ID"
//	@Success	200			{object}	object
//	@Failure	404			{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/subjects/{subjectID} [get]
func (h *Handler) GetSubject(w http.ResponseWriter, r *http.Request) {
	subjectID := chi.URLParam(r, "subjectID")
	subject, err := h.SubjectService.GetSubject(subjectID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Subject not found")
		return
	}
	h.writeJSON(w, http.StatusOK, subject)
}

// CreateSubject handles POST /v1/subjects.
//
//	@Summary	Create a new subject
//	@Tags		subjects
//	@Accept		json
//	@Produce	json
//	@Param		request	body		SubjectCreateRequest	true	"Create subject request"
//	@Success	201		{object}	StatusResponse
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/subjects [post]
func (h *Handler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	var req SubjectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.SubjectName == "" {
		h.writeError(w, http.StatusBadRequest, "subject_name is required")
		return
	}
	id, err := h.SubjectService.CreateSubject(req.SubjectName)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, StatusResponse{ID: id, Name: req.SubjectName})
}

// UpdateSubject handles PUT /v1/subjects/{subjectID}.
//
//	@Summary	Update a subject
//	@Tags		subjects
//	@Accept		json
//	@Produce	json
//	@Param		subjectID	path		string					true	"Subject ID"
//	@Param		request		body		SubjectUpdateRequest	true	"Update subject request"
//	@Success	200			{object}	StatusResponse
//	@Failure	400			{object}	ErrorResponse
//	@Failure	404			{object}	ErrorResponse
//	@Failure	500			{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/subjects/{subjectID} [put]
func (h *Handler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
	subjectID := chi.URLParam(r, "subjectID")
	var req SubjectUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	subject, err := h.SubjectService.GetSubject(subjectID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Subject not found")
		return
	}
	subject.SubjectName = &req.SubjectName
	if err := h.SubjectService.SaveSubject(subject); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, StatusResponse{ID: subjectID, Status: "updated"})
}

// DeleteSubject handles DELETE /v1/subjects/{subjectID}.
//
//	@Summary	Delete a subject
//	@Tags		subjects
//	@Produce	json
//	@Param		subjectID	path		string	true	"Subject ID"
//	@Success	200			{object}	StatusResponse
//	@Failure	404			{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/subjects/{subjectID} [delete]
func (h *Handler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
	subjectID := chi.URLParam(r, "subjectID")
	if err := h.SubjectService.DeleteSubject(subjectID); err != nil {
		h.writeError(w, http.StatusNotFound, "Subject not found")
		return
	}
	h.writeJSON(w, http.StatusOK, StatusResponse{ID: subjectID, Status: "deleted"})
}
