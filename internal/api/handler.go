package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sovraai/subject-data/internal/service"
)

// Handler holds dependencies for the HTTP handlers.
type Handler struct {
	Logger         *slog.Logger
	SubjectService *service.SubjectService
	RecordService  *service.RecordService
}

// New constructs a Handler.
func New(logger *slog.Logger, subjectSvc *service.SubjectService, recordSvc *service.RecordService) *Handler {
	return &Handler{
		Logger:         logger,
		SubjectService: subjectSvc,
		RecordService:  recordSvc,
	}
}

// Healthz is the liveness probe.
//
//	@Summary	Liveness probe
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/healthz [get]
func (h *Handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Readyz is the readiness probe.
//
//	@Summary	Readiness probe
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/readyz [get]
func (h *Handler) Readyz(w http.ResponseWriter, _ *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.Logger.Error("write json", "err", err)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, ErrorResponse{
		Error: ErrorDetail{
			Code:    http.StatusText(status),
			Message: message,
		},
	})
}
