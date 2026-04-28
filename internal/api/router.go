package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	datamodel "github.com/sovraai/subject-data/data-model"
	"github.com/sovraai/subject-data/internal/middleware"
	"github.com/sovraai/subject-data/internal/service"
	"github.com/sovraai/subject-data/swagger"
	"github.com/sovraai/subject-data/swagger/schemabrowser"
)

// RouterDeps holds all dependencies needed by route handlers.
type RouterDeps struct {
	Logger         *slog.Logger
	AuthTokens     []string
	CORSOrigins    string
	SubjectService *service.SubjectService
	RecordService  *service.RecordService
}

// NewRouter creates and configures the Chi router with all API routes.
func NewRouter(deps RouterDeps) *chi.Mux {
	r := chi.NewRouter()

	// CORS — dev-only convenience so a local frontend dev server can hit
	// subject-data directly from the browser. This service is
	// server-to-server in production; CORS is irrelevant behind an ALB.
	if deps.CORSOrigins != "" {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{deps.CORSOrigins},
			AllowedMethods: []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
			MaxAge:         300,
		}))
		deps.Logger.Info("CORS enabled", "origins", deps.CORSOrigins)
	}

	h := New(deps.Logger, deps.SubjectService, deps.RecordService)

	// Health endpoints — no auth.
	r.Get("/healthz", h.Healthz)
	r.Get("/readyz", h.Readyz)

	// API docs — Swagger UI for REST endpoints, no auth.
	swagger.MountDocs(r, "/docs/doc.json")

	// Schema browser — data model type documentation, no auth.
	r.Route("/schema-browser", func(r chi.Router) {
		schemabrowser.Mount(r, schemabrowser.Config{
			Title:   "Subject Data — Schema Browser",
			Version: "0.1.0",
			Models:  subjectDataModels(),
		})
	})

	// Authenticated endpoints — Bearer token from AUTH_TOKENS allowlist.
	// If AUTH_TOKENS is empty, auth is disabled (local dev without a token).
	r.Group(func(r chi.Router) {
		r.Use(middleware.BearerAuth(deps.AuthTokens))

		r.Route("/v1", func(r chi.Router) {
			r.Get("/subjects", h.ListSubjects)
			r.Post("/subjects", h.CreateSubject)
			r.Get("/subjects/{subjectID}", h.GetSubject)
			r.Put("/subjects/{subjectID}", h.UpdateSubject)
			r.Delete("/subjects/{subjectID}", h.DeleteSubject)

			r.Get("/subjects/{subjectID}/records", h.ListRecordsBySubject)

			r.Get("/records", h.ListRecords)
			r.Post("/records", h.CreateRecord)
			r.Get("/records/{recordID}", h.GetRecord)
			r.Put("/records/{recordID}", h.UpdateRecord)
			r.Delete("/records/{recordID}", h.DeleteRecord)
		})
	})

	return r
}

// subjectDataModels returns the model entries for the schema browser,
// including all data-model types with their examples.
func subjectDataModels() []schemabrowser.ModelEntry {
	examples := datamodel.Examples()
	models := []struct {
		Name     string
		Tag      string
		Instance any
	}{
		// Subject
		{"Subject", "Subject", datamodel.Subject{}},
		// Records
		{"Record", "Records", datamodel.Record{}},
		{"RecordMetadata", "Records", datamodel.RecordMetadata{}},
		{"DataSource", "Records", datamodel.DataSource{}},
		{"PlatformData", "Records", datamodel.PlatformData{}},
		{"MediaAnalysis", "Records", datamodel.MediaAnalysis{}},
		{"ImageEntry", "Records", datamodel.ImageEntry{}},
		{"VideoEntry", "Records", datamodel.VideoEntry{}},

		// CE Features & Ratings
		{"CEFeatureDefinition", "CE Features", datamodel.CEFeatureDefinition{}},
		{"CEFeatureDefinitionEnriched", "CE Features", datamodel.CEFeatureDefinitionEnriched{}},
		{"SubjectCEFeature", "CE Features", datamodel.SubjectCEFeature{}},
		{"Evidence", "CE Features", datamodel.Evidence{}},
		{"EvidenceEnriched", "CE Features", datamodel.EvidenceEnriched{}},
		{"RatingBenchmark", "CE Features", datamodel.RatingBenchmark{}},
		{"SubjectRatings", "CE Features", datamodel.SubjectRatings{}},
		{"SubjectRatingsEnriched", "CE Features", datamodel.SubjectRatingsEnriched{}},

		// Joined Data
		{"JoinedSubjectData", "Joined Data", datamodel.JoinedSubjectData{}},
		{"JoinedDatasetEnriched", "Joined Data", datamodel.JoinedDatasetEnriched{}},
	}

	entries := make([]schemabrowser.ModelEntry, len(models))
	for i, m := range models {
		entries[i] = schemabrowser.ModelEntry{
			Name:     m.Name,
			Tag:      m.Tag,
			Instance: m.Instance,
			Example:  examples[m.Name],
		}
	}
	return entries
}
