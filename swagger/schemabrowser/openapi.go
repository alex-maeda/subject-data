package schemabrowser

import (
	"encoding/json"
	"strings"

	"github.com/invopop/jsonschema"
)

// GenerateSpec builds an OpenAPI 3.0.0 JSON spec from the given Config.
// The spec includes a schema component and GET path for each ModelEntry.
func GenerateSpec(cfg Config) []byte {
	r := jsonschema.Reflector{}

	tagSet := map[string]bool{}
	var tagList []map[string]string

	components := map[string]any{}

	// Reflect all models and collect schemas.
	for _, m := range cfg.Models {
		schema := r.Reflect(m.Instance)
		schemaBytes, _ := json.Marshal(schema)
		var schemaMap map[string]any
		json.Unmarshal(schemaBytes, &schemaMap) //nolint:errcheck // schema is valid JSON

		// Extract $defs (nested type definitions) into components.
		defs, hasDefs := schemaMap["$defs"].(map[string]any)
		if hasDefs {
			for name, def := range defs {
				if _, exists := components[name]; !exists {
					defMap := rewriteRefs(def).(map[string]any)
					if m.Example != nil {
						// Only set example on the top-level model, not nested defs
					}
					components[name] = defMap
				}
			}
		}

		// The top-level schema may be a $ref to $defs. Resolve it.
		delete(schemaMap, "$defs")
		delete(schemaMap, "$schema")
		delete(schemaMap, "$id")

		topLevel := schemaMap
		if ref, ok := schemaMap["$ref"].(string); ok && hasDefs {
			refName := strings.TrimPrefix(ref, "#/$defs/")
			if resolved, ok := defs[refName]; ok {
				topLevel = rewriteRefs(resolved).(map[string]any)
			}
		} else {
			topLevel = rewriteRefs(schemaMap).(map[string]any)
		}

		if m.Example != nil {
			topLevel["example"] = m.Example
		}
		components[m.Name] = topLevel
	}

	// Build paths — each model gets a GET endpoint.
	paths := map[string]any{}
	for _, m := range cfg.Models {
		paths["/"+m.Name] = map[string]any{
			"get": map[string]any{
				"summary":     m.Name,
				"operationId": "get" + m.Name,
				"tags":        []string{m.Tag},
				"responses": map[string]any{
					"200": map[string]any{
						"description": m.Name + " schema",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]string{
									"$ref": "#/components/schemas/" + m.Name,
								},
							},
						},
					},
				},
			},
		}

		if !tagSet[m.Tag] {
			tagSet[m.Tag] = true
			tagList = append(tagList, map[string]string{"name": m.Tag})
		}
	}

	spec := map[string]any{
		"openapi": "3.0.0",
		"info": map[string]any{
			"title":   cfg.Title,
			"version": cfg.Version,
		},
		"paths": paths,
		"tags":  tagList,
		"components": map[string]any{
			"schemas": components,
		},
	}

	data, _ := json.MarshalIndent(spec, "", "  ")
	return data
}

// rewriteRefs converts JSON Schema $defs references to OpenAPI components/schemas references.
func rewriteRefs(v any) any {
	switch val := v.(type) {
	case map[string]any:
		result := make(map[string]any, len(val))
		for k, child := range val {
			if k == "$ref" {
				if s, ok := child.(string); ok {
					result[k] = strings.Replace(s, "#/$defs/", "#/components/schemas/", 1)
				} else {
					result[k] = child
				}
			} else {
				result[k] = rewriteRefs(child)
			}
		}
		return result
	case []any:
		result := make([]any, len(val))
		for i, child := range val {
			result[i] = rewriteRefs(child)
		}
		return result
	default:
		return v
	}
}
