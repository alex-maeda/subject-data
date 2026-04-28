// Package schemabrowser generates and serves an OpenAPI spec for Go types,
// providing a browsable schema documentation UI via Swagger UI (CDN).
//
// Usage:
//
//	r.Route("/schema-browser", func(r chi.Router) {
//	    schemabrowser.Mount(r, schemabrowser.Config{
//	        Title: "My Data Models",
//	        Models: []schemabrowser.ModelEntry{
//	            {Name: "User", Tag: "Core", Instance: User{}, Example: exampleUser},
//	        },
//	    })
//	})
package schemabrowser

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ModelEntry describes a Go type to include in the schema browser.
type ModelEntry struct {
	Name     string // Schema name and URL path segment, e.g. "Subject"
	Tag      string // OpenAPI tag for grouping in the UI, e.g. "Core"
	Instance any    // Zero-value instance for jsonschema reflection
	Example  any    // Optional: populated instance for GET /{Name} endpoint
}

// Config holds configuration for the schema browser.
type Config struct {
	Models  []ModelEntry
	Title   string // OpenAPI info.title (default: "Schema Browser")
	Version string // OpenAPI info.version (default: "1.0.0")
}

const swaggerUIHTML = `<!DOCTYPE html>
<html>
<head>
  <title>%s</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: 'openapi.json',
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
      layout: 'BaseLayout',
    });
  </script>
</body>
</html>`

// Mount registers schema browser routes on the given chi.Router:
//   - GET /docs           — Swagger UI pointing at ./openapi.json
//   - GET /openapi.json   — generated OpenAPI 3.0.0 spec
//   - GET /{ModelName}    — example JSON (if Example was provided on ModelEntry)
func Mount(r chi.Router, cfg Config) {
	if cfg.Title == "" {
		cfg.Title = "Schema Browser"
	}
	if cfg.Version == "" {
		cfg.Version = "1.0.0"
	}

	specJSON := GenerateSpec(cfg)

	r.Get("/docs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//nolint:errcheck // best-effort write
		w.Write([]byte(fmt.Sprintf(swaggerUIHTML, cfg.Title)))
	})

	r.Get("/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//nolint:errcheck // best-effort write
		w.Write(specJSON)
	})

	for _, m := range cfg.Models {
		example := m.Example
		if example == nil {
			continue
		}
		r.Get("/"+m.Name, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			data, _ := json.Marshal(example)
			//nolint:errcheck // best-effort write
			w.Write(data)
		})
	}
}
