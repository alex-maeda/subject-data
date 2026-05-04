package datamodel

import "encoding/json"

// AttributeType is the closed enum of supported Attribute discriminator values.

// String values use kebab-case to match the BFF's card slugs and the
// `?attribute_names=relationship,social-media-footprint,...` query
// parameter — no translation layer needed.
type AttributeType string

const (
	// Page header — atomic facts populated by SO and DE
	AttributeTypeName   AttributeType = "name"
	AttributeTypeAge    AttributeType = "age"
	AttributeTypeGender AttributeType = "gender"
	AttributeTypePhone  AttributeType = "phone"
	AttributeTypeEmail  AttributeType = "email"
	AttributeTypeAvatar AttributeType = "avatar"

	// Cards — SO-written (sync, beta-scope)
	AttributeTypeRelationship         AttributeType = "relationship"
	AttributeTypeSocialMediaFootprint AttributeType = "social-media-footprint"
	AttributeTypeGeographicFootprint  AttributeType = "geographic-footprint"

	// Cards — DE-written (async)
	AttributeTypePublicRecords AttributeType = "public-records"
	AttributeTypeTimelines     AttributeType = "timelines"
	AttributeTypeInTheNews     AttributeType = "in-the-news"
	AttributeTypeDataLeaks     AttributeType = "data-leaks"

	// Cards — AE-written (post-DE)
	AttributeTypeBehavioralAnalysis AttributeType = "behavioral-analysis"
	AttributeTypeSummary            AttributeType = "summary"
)

// Attribute is a typed, schema'd, confidence-bearing value attached to a Subject.

// Different attribute types (name, phone, public-records, ...) have different
// payload shapes, but every Attribute shares this envelope. The payload is
// stored as JSON and validated server-side at write time against the registered
// schema for the (Type, SchemaVersion) pair.

// Confidence and EvidenceSummary live on the envelope, not in the payload.
// Per-payload domain-specific confidences (per-account match_confidence,
// per-record severity, etc.) are separate concepts and stay in payloads.
type Attribute struct {
	Identifiable

	SubjectID       string                 `json:"subject_id" jsonschema_description:"Identifier of the Subject this Attribute is attached to"`
	Type            AttributeType          `json:"type" jsonschema_description:"Kind of Attribute (matches BFF card slug)"`
	SchemaVersion   string                 `json:"schema_version" jsonschema_description:"Version of the per-type payload schema (e.g. v1)"`
	Payload         json.RawMessage        `json:"payload" jsonschema_description:"Type-specific structured value, validated against the schema for (Type, SchemaVersion)"`
	Confidence      float64                `json:"confidence" jsonschema_description:"Confidence in this Attribute's value, 0.0 to 1.0"`
	EvidenceSummary string                 `json:"evidence_summary" jsonschema_description:"Human-readable explanation of how this Attribute was derived"`
	EvidenceRefs    []AttributeEvidenceRef `json:"evidence_refs" jsonschema_description:"Per-record provenance pointers"`
	CreatedAt       DateTime               `json:"created_at"`
	UpdatedAt       DateTime               `json:"updated_at"`
}

// NewAttribute creates an Attribute with initialized slices.
func NewAttribute() Attribute {
	return Attribute{
		EvidenceRefs: []AttributeEvidenceRef{},
	}
}

// AttributeEvidenceRef is one column-level provenance pointer into a contributing
// Record. Mirrors IDS's evidence triple (source_field, normalized_attribute,
// normalized_value) but adds the cross-record dimension via record_id and drops
// the value — values live on the Record and are looked up by record_id, not
// duplicated here.
type AttributeEvidenceRef struct {
	RecordID            string `json:"record_id" jsonschema_description:"Identifier of the contributing Record"`
	SourceField         string `json:"source_field" jsonschema_description:"Name of the column in that Record that produced the fact"`
	NormalizedAttribute string `json:"normalized_attribute,omitempty" jsonschema_description:"Name of the derived attribute (mirrors IDS's evidence contract); often redundant with Attribute.Type but preserved for traceability"`
}

// Sanitizer is implemented by Attribute payloads that contain raw PII fields
// requiring redaction before API responses. Read handlers decode the payload
// into its concrete type and call Sanitize() to null out raw fields, then
// re-encode for the response. Internal callers needing raw values read
// Records directly, not Attributes.
type Sanitizer interface {
	Sanitize()
}
