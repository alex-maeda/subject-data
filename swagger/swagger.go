// Package swagger provides reusable Swagger UI hosting for Go HTTP services.
//
// Use MountDocs to serve the Swagger UI for any OpenAPI spec, or use the
// schemabrowser subpackage to auto-generate a browsable schema from Go types.
package swagger

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// MountDocs mounts Swagger UI at /docs on the given router.
// specURL is the URL the UI will fetch the OpenAPI spec from
// (e.g. "/docs/doc.json" for swaggo-generated specs).
func MountDocs(r chi.Router, specURL string) {
	r.Get("/docs", http.RedirectHandler("/docs/index.html", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL(specURL),
	))
}
