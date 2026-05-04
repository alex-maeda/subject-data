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

		// Attributes — envelope
		{"Attribute", "Attributes", datamodel.Attribute{}},
		{"AttributeEvidenceRef", "Attributes", datamodel.AttributeEvidenceRef{}},

		// Attributes — header payloads (Tier 1)
		{"AttributeNamePayloadV1", "Attribute Payloads", datamodel.AttributeNamePayloadV1{}},
		{"AttributeAgePayloadV1", "Attribute Payloads", datamodel.AttributeAgePayloadV1{}},
		{"AttributeGenderPayloadV1", "Attribute Payloads", datamodel.AttributeGenderPayloadV1{}},
		{"AttributePhonePayloadV1", "Attribute Payloads", datamodel.AttributePhonePayloadV1{}},
		{"AttributeEmailPayloadV1", "Attribute Payloads", datamodel.AttributeEmailPayloadV1{}},
		{"AttributeAvatarPayloadV1", "Attribute Payloads", datamodel.AttributeAvatarPayloadV1{}},

		// Attributes — SO-written card payloads (Tier 2)
		{"AttributeRelationshipPayloadV1", "Attribute Payloads", datamodel.AttributeRelationshipPayloadV1{}},
		{"AttributeSocialMediaFootprintPayloadV1", "Attribute Payloads", datamodel.AttributeSocialMediaFootprintPayloadV1{}},
		{"SocialAccount", "Attribute Payloads", datamodel.SocialAccount{}},
		{"AttributeGeographicFootprintPayloadV1", "Attribute Payloads", datamodel.AttributeGeographicFootprintPayloadV1{}},
		{"GeoCoordinates", "Attribute Payloads", datamodel.GeoCoordinates{}},
		{"CityLocation", "Attribute Payloads", datamodel.CityLocation{}},
		{"LocationFrequency", "Attribute Payloads", datamodel.LocationFrequency{}},
		{"GeographicLocation", "Attribute Payloads", datamodel.GeographicLocation{}},
		{"GeographicDataPoint", "Attribute Payloads", datamodel.GeographicDataPoint{}},

		// Attributes — DE-written card payloads (Tier 2)
		{"AttributePublicRecordsPayloadV1", "Attribute Payloads", datamodel.AttributePublicRecordsPayloadV1{}},
		{"PublicRecordGroup", "Attribute Payloads", datamodel.PublicRecordGroup{}},
		{"PublicRecord", "Attribute Payloads", datamodel.PublicRecord{}},
		{"PublicRecordParty", "Attribute Payloads", datamodel.PublicRecordParty{}},
		{"AttributeTimelinesPayloadV1", "Attribute Payloads", datamodel.AttributeTimelinesPayloadV1{}},
		{"TimelineEvent", "Attribute Payloads", datamodel.TimelineEvent{}},

		// Attributes — stubs (Tier 3 — post-beta / post-GA)
		{"AttributeInTheNewsPayloadV1", "Attribute Payloads", datamodel.AttributeInTheNewsPayloadV1{}},
		{"NewsArticle", "Attribute Payloads", datamodel.NewsArticle{}},
		{"AttributeDataLeaksPayloadV1", "Attribute Payloads", datamodel.AttributeDataLeaksPayloadV1{}},
		{"DataLeak", "Attribute Payloads", datamodel.DataLeak{}},
		{"DataLeakSource", "Attribute Payloads", datamodel.DataLeakSource{}},
		{"AttributeBehavioralAnalysisPayloadV1", "Attribute Payloads", datamodel.AttributeBehavioralAnalysisPayloadV1{}},
		{"AttributeSummaryPayloadV1", "Attribute Payloads", datamodel.AttributeSummaryPayloadV1{}},
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
