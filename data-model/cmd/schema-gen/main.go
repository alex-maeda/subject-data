// Command schema-gen generates JSON Schema files for all data model types.
//
// Usage:
//
//	go run ./cmd/schema-gen <output-dir>
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"

	datamodel "github.com/sovraai/subject-data/data-model"
)

type modelEntry struct {
	Name     string
	Instance interface{}
}

func models() []modelEntry {
	return []modelEntry{
		{"CEFeatureDefinition", datamodel.CEFeatureDefinition{}},
		{"CEFeatureDefinitionEnriched", datamodel.CEFeatureDefinitionEnriched{}},
		{"SubjectCEFeature", datamodel.SubjectCEFeature{}},
		{"EvidenceEnriched", datamodel.EvidenceEnriched{}},
		{"Subject", datamodel.Subject{}},
		{"Record", datamodel.Record{}},
		{"SubjectRatings", datamodel.SubjectRatings{}},
		{"SubjectRatingsEnriched", datamodel.SubjectRatingsEnriched{}},
		{"JoinedSubjectData", datamodel.JoinedSubjectData{}},
		{"JoinedDatasetEnriched", datamodel.JoinedDatasetEnriched{}},
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <output-dir>\n", os.Args[0])
		os.Exit(1)
	}
	outDir := os.Args[1]
	if err := os.MkdirAll(outDir, 0o750); err != nil { //nolint:gosec // CLI tool, user-supplied output dir is intentional
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	r := jsonschema.Reflector{}

	for _, m := range models() {
		schema := r.Reflect(m.Instance)
		data, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling %s: %v\n", m.Name, err)
			os.Exit(1)
		}
		path := filepath.Join(outDir, m.Name+".json")
		if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil { //nolint:gosec // CLI tool, output path derived from user-supplied dir
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("  %s.json\n", m.Name)
	}
	fmt.Println("Done.")
}
