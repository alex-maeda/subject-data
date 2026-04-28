package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	datamodel "github.com/sovraai/subject-data/data-model"
)

// ListRecordsBySubject handles GET /v1/subjects/{subjectID}/records.
//
//	@Summary	List records for a subject
//	@Tags		records
//	@Produce	json
//	@Param		subjectID	path		string	true	"Subject ID"
//	@Success	200			{array}		object
//	@Failure	500			{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/subjects/{subjectID}/records [get]
func (h *Handler) ListRecordsBySubject(w http.ResponseWriter, r *http.Request) {
	subjectID := chi.URLParam(r, "subjectID")
	data, err := h.RecordService.ListBySubject(subjectID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, data)
}

// ListRecords handles GET /v1/records.
//
//	@Summary	List all records
//	@Tags		records
//	@Produce	json
//	@Param		subject_id	query		string	false	"Filter by subject ID"
//	@Success	200			{array}		object
//	@Failure	500			{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/records [get]
func (h *Handler) ListRecords(w http.ResponseWriter, r *http.Request) {
	subjectID := r.URL.Query().Get("subject_id")
	var data []datamodel.Record
	var err error
	if subjectID != "" {
		data, err = h.RecordService.ListBySubject(subjectID)
	} else {
		data, err = h.RecordService.ListAll()
	}
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, data)
}

// CreateRecord handles POST /v1/records.
//
//	@Summary	Create a new record
//	@Tags		records
//	@Accept		json
//	@Produce	json
//	@Param		request	body		object	true	"Record data"
//	@Success	201		{object}	StatusResponse
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/records [post]
func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	var obj datamodel.Record
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	id, err := h.RecordService.Create(&obj)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, StatusResponse{ID: id})
}

// GetRecord handles GET /v1/records/{recordID}.
//
//	@Summary	Get record by ID
//	@Tags		records
//	@Produce	json
//	@Param		recordID	path		string	true	"Record ID"
//	@Success	200			{object}	object
//	@Failure	404			{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/records/{recordID} [get]
func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) {
	recordID := chi.URLParam(r, "recordID")
	obj, err := h.RecordService.Get(recordID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Record not found")
		return
	}
	h.writeJSON(w, http.StatusOK, obj)
}

// UpdateRecord handles PUT /v1/records/{recordID}.
//
//	@Summary	Update a record
//	@Tags		records
//	@Accept		json
//	@Produce	json
//	@Param		recordID	path		string	true	"Record ID"
//	@Param		request		body		object	true	"Updated record data"
//	@Success	200			{object}	StatusResponse
//	@Failure	400			{object}	ErrorResponse
//	@Failure	404			{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/records/{recordID} [put]
func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	recordID := chi.URLParam(r, "recordID")
	var obj datamodel.Record
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.RecordService.Update(recordID, &obj); err != nil {
		h.writeError(w, http.StatusNotFound, "Record not found")
		return
	}
	h.writeJSON(w, http.StatusOK, StatusResponse{ID: recordID, Status: "updated"})
}

// DeleteRecord handles DELETE /v1/records/{recordID}.
//
//	@Summary	Delete a record
//	@Tags		records
//	@Produce	json
//	@Param		recordID	path		string	true	"Record ID"
//	@Success	200			{object}	StatusResponse
//	@Failure	404			{object}	ErrorResponse
//	@Security	BearerAuth
//	@Router		/v1/records/{recordID} [delete]
func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	recordID := chi.URLParam(r, "recordID")
	if err := h.RecordService.Delete(recordID); err != nil {
		h.writeError(w, http.StatusNotFound, "Record not found")
		return
	}
	h.writeJSON(w, http.StatusOK, StatusResponse{ID: recordID, Status: "deleted"})
}
