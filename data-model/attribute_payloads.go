package datamodel

// =============================================================================
// Attribute payloads — v1 schemas for each AttributeType.
//
// Conventions:
//   - Each payload type is named AttributeXxxPayloadV1 where Xxx matches the
//     AttributeType const (PascalCase).
//   - Field names mirror BFF TypeScript types (snake_case in Go/JSON;
//     camelCase in TS — BFF maps).
//   - Payloads with raw PII fields implement Sanitizer to redact those fields
//     before API responses.
//   - Domain-specific confidences (per-account match_confidence, per-leak
//     severity, etc.) stay in payloads. Card-level confidence lives on the
//     Attribute envelope, not duplicated here.
//
// Tier breakdown (for review reference):
//   Tier 1 — Header atomic facts. Real schemas, beta-blocking.
//   Tier 2 — SO/DE-written cards. Real schemas for beta scope.
//   Tier 3 — Stubs for post-beta / post-GA cards.
// =============================================================================

// -----------------------------------------------------------------------------
// Tier 1 — Header attributes
// -----------------------------------------------------------------------------

// AttributeNamePayloadV1 — payload for type=name.
type AttributeNamePayloadV1 struct {
	First       string `json:"first,omitempty"`
	Middle      string `json:"middle,omitempty"`
	Last        string `json:"last,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	Suffix      string `json:"suffix,omitempty"`
	DisplayName string `json:"display_name,omitempty" jsonschema_description:"Optional precomputed display variant; BFF can fall back to first + last"`
}

// AttributeAgePayloadV1 — payload for type=age. Stores year_of_birth;
// BFF computes current age.
type AttributeAgePayloadV1 struct {
	YearOfBirth int    `json:"year_of_birth" jsonschema_description:"Year of birth; BFF computes current age from this"`
	BirthDate   string `json:"birth_date,omitempty" jsonschema_description:"Optional ISO 8601 YYYY-MM-DD when known precisely"`
}

// Gender is the closed enum for AttributeGenderPayloadV1.Value.
type Gender string

const (
	GenderMale      Gender = "male"
	GenderFemale    Gender = "female"
	GenderNonBinary Gender = "non_binary"
	GenderOther     Gender = "other"
	GenderUnknown   Gender = "unknown"
)

// AttributeGenderPayloadV1 — payload for type=gender.
type AttributeGenderPayloadV1 struct {
	Value Gender `json:"value"`
}

// AttributePhonePayloadV1 — payload for type=phone.
//
// PII rule: E164 is the raw value; Obfuscated is the BFF-rendering form.
// Sanitize() is called by the API layer before responses to strip E164.
// Internal callers needing raw E164 read Records directly.
type AttributePhonePayloadV1 struct {
	E164       string `json:"e164,omitempty" jsonschema_description:"Raw E.164 phone; stripped from API responses"`
	Obfuscated string `json:"obfuscated" jsonschema_description:"BFF-rendering form, e.g. '+1 617 ***-***3'"`
	Country    string `json:"country,omitempty" jsonschema_description:"ISO 3166-1 alpha-2 country code"`
}

// Sanitize zeroes the raw E164 field. Called before API responses.
func (p *AttributePhonePayloadV1) Sanitize() {
	p.E164 = ""
}

// AttributeEmailPayloadV1 — payload for type=email.

// PII rule: Address is the raw value; Obfuscated is the BFF-rendering form.
// Sanitize() is called by the API layer before responses to strip Address.
type AttributeEmailPayloadV1 struct {
	Address    string `json:"address,omitempty" jsonschema_description:"Raw email address; stripped from API responses"`
	Obfuscated string `json:"obfuscated" jsonschema_description:"BFF-rendering form, e.g. 'sm********@gmail.com'"`
}

// Sanitize zeroes the raw Address field. Called before API responses.
func (p *AttributeEmailPayloadV1) Sanitize() {
	p.Address = ""
}

// AttributeAvatarPayloadV1 — payload for type=avatar.
type AttributeAvatarPayloadV1 struct {
	URL string `json:"url"`
}

// -----------------------------------------------------------------------------
// Tier 2 — SO-written cards (beta-scope)
// -----------------------------------------------------------------------------

// RelationshipStatus is the closed enum for AttributeRelationshipPayloadV1.Status.
//
// v1 minimum: unknown / married / not_married. Future values (divorced,
// widowed, engaged, in_relationship) deferred — Request confirmation of the
// stakeholders on the post-beta taxonomy.
type RelationshipStatus string

const (
	RelationshipStatusUnknown    RelationshipStatus = "unknown"
	RelationshipStatusMarried    RelationshipStatus = "married"
	RelationshipStatusNotMarried RelationshipStatus = "not_married"
)

// AttributeRelationshipPayloadV1 — payload for type=relationship.
//
// Confidence and evidenceSummary live on the Attribute envelope; the
// payload doesn't duplicate them.
type AttributeRelationshipPayloadV1 struct {
	Status RelationshipStatus `json:"status"`
}

// SocialPlatform is the closed enum for SocialAccount.Platform.
type SocialPlatform string

const (
	SocialPlatformLinkedIn  SocialPlatform = "linkedin"
	SocialPlatformInstagram SocialPlatform = "instagram"
	SocialPlatformTwitter   SocialPlatform = "twitter"
	SocialPlatformFacebook  SocialPlatform = "facebook"
	SocialPlatformTikTok    SocialPlatform = "tiktok"
	SocialPlatformBluesky   SocialPlatform = "bluesky"
	SocialPlatformReddit    SocialPlatform = "reddit"
	SocialPlatformThreads   SocialPlatform = "threads"
)

// SocialMatchConfidence is the binary per-account match confidence.
type SocialMatchConfidence string

const (
	SocialMatchVerified SocialMatchConfidence = "verified"
	SocialMatchInferred SocialMatchConfidence = "inferred"
)

// SocialAccount is one row in social-media-footprint.accounts.
type SocialAccount struct {
	Platform        SocialPlatform        `json:"platform"`
	Handle          string                `json:"handle"`
	ProfileURL      string                `json:"profile_url"`
	MatchConfidence SocialMatchConfidence `json:"match_confidence"`
}

// AttributeSocialMediaFootprintPayloadV1 — payload for type=social-media-footprint.
//
// Beta scope (Tier 1): handle + profileURL + matchConfidence only. No follower
// counts, no profile photos, no posts.
type AttributeSocialMediaFootprintPayloadV1 struct {
	Accounts []SocialAccount `json:"accounts"`
}

// GeoCoordinates is a lat/lng pair. Optional in v1 — geocoding is post-beta polish.
type GeoCoordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// CityLocation is the BFF "header / pill" form — city + state, optional coords.
type CityLocation struct {
	City        string          `json:"city,omitempty"`
	State       string          `json:"state,omitempty"`
	Coordinates *GeoCoordinates `json:"coordinates,omitempty"`
}

// LocationFrequency is one entry in geographic-footprint.location_frequency.
type LocationFrequency struct {
	City        string          `json:"city,omitempty"`
	State       string          `json:"state,omitempty"`
	Count       int             `json:"count"`
	Coordinates *GeoCoordinates `json:"coordinates,omitempty"`
}

// GeographicLocation is the fuller location form used inside data points,
// including the raw address string.
//
// PII rule: Address is raw and stripped from API responses by the
// AttributeGeographicFootprintPayloadV1.Sanitize implementation.
type GeographicLocation struct {
	City        string          `json:"city,omitempty"`
	State       string          `json:"state,omitempty"`
	Address     string          `json:"address,omitempty" jsonschema_description:"Raw street address; stripped from API responses"`
	Coordinates *GeoCoordinates `json:"coordinates,omitempty"`
}

// GeographicDataPoint is one entry in geographic-footprint.data_points.
type GeographicDataPoint struct {
	Location    GeographicLocation `json:"location"`
	Date        string             `json:"date,omitempty"`
	SourceType  string             `json:"source_type,omitempty"`
	SourceRef   string             `json:"source_ref,omitempty"`
	Description string             `json:"description,omitempty"`
}

// AttributeGeographicFootprintPayloadV1 — payload for type=geographic-footprint.

// Beta scope: latest_location, location_frequency, data_points from public
// records. Map rendering with Coordinates is post-beta polish.
type AttributeGeographicFootprintPayloadV1 struct {
	LatestLocation    *CityLocation         `json:"latest_location,omitempty"`
	LocationFrequency []LocationFrequency   `json:"location_frequency"`
	DataPoints        []GeographicDataPoint `json:"data_points"`
}

// Sanitize strips raw addresses from data_points before API responses.
func (p *AttributeGeographicFootprintPayloadV1) Sanitize() {
	for i := range p.DataPoints {
		p.DataPoints[i].Location.Address = ""
	}
}

// -----------------------------------------------------------------------------
// Tier 2 — DE-written cards (beta-scope)
// -----------------------------------------------------------------------------

// PublicRecordType is the closed enum for PublicRecordGroup.Type.
type PublicRecordType string

const (
	PublicRecordTypeCivil    PublicRecordType = "civil"
	PublicRecordTypeCriminal PublicRecordType = "criminal"
	PublicRecordTypeMarriage PublicRecordType = "marriage"
	// Future: bankruptcy, mortgage_deed, family_court, etc. — extend as DE adds sources.
)

// PublicRecordParty is one named participant in a public record (co-defendant,
// counter-plaintiff, witness, etc.). Names are PII and stripped on response.
type PublicRecordParty struct {
	Name string `json:"name" jsonschema_description:"Party name; PII, stripped from API responses"`
	Role string `json:"role,omitempty"`
}

// PublicRecord is one record in a public_records group.
//
// Required-ish fields per BFF's spec: title, date, location, description,
// role, classification, caseType, disposition.
// Nice-to-have: caseNumber, court, sourceUrl, parties, sentence/severity (criminal-only).
type PublicRecord struct {
	Title          string              `json:"title"`
	Date           string              `json:"date,omitempty"`
	Location       string              `json:"location,omitempty"`
	Description    string              `json:"description,omitempty"`
	Role           string              `json:"role,omitempty" jsonschema_description:"Subject's role: plaintiff, defendant, petitioner, etc."`
	Classification string              `json:"classification,omitempty" jsonschema_description:"e.g. settled, dismissed, convicted, acquitted"`
	CaseType       string              `json:"case_type,omitempty"`
	Disposition    string              `json:"disposition,omitempty"`
	Parties        []PublicRecordParty `json:"parties,omitempty"`
	CaseNumber     string              `json:"case_number,omitempty"`
	Court          string              `json:"court,omitempty"`
	SourceURL      string              `json:"source_url,omitempty"`
	Sentence       string              `json:"sentence,omitempty" jsonschema_description:"Criminal records only"`
	Severity       string              `json:"severity,omitempty" jsonschema_description:"Criminal records only"`
}

// PublicRecordGroup groups records by type (civil / criminal / marriage / ...).
type PublicRecordGroup struct {
	Type    PublicRecordType `json:"type"`
	Records []PublicRecord   `json:"records"`
}

// AttributePublicRecordsPayloadV1 — payload for type=public-records.
type AttributePublicRecordsPayloadV1 struct {
	Groups []PublicRecordGroup `json:"groups"`
}

// Sanitize strips party names from records before API responses.
// Other PII (case numbers, descriptions) is left intact for v1; if it becomes
// a concern, extend this method.
func (p *AttributePublicRecordsPayloadV1) Sanitize() {
	for gi := range p.Groups {
		for ri := range p.Groups[gi].Records {
			rec := &p.Groups[gi].Records[ri]
			for pi := range rec.Parties {
				rec.Parties[pi].Name = ""
			}
		}
	}
}

// TimelineCategory is the closed enum for TimelineEvent.Category.
type TimelineCategory string

const (
	TimelineCategoryLocation     TimelineCategory = "location"
	TimelineCategoryProfessional TimelineCategory = "professional"
	TimelineCategoryEducation    TimelineCategory = "education"
	TimelineCategoryPersonal     TimelineCategory = "personal"
)

// TimelineEvent is one event in the timelines card. Beta = location lane only;
// other categories are populated post-beta when DE has the sources.
type TimelineEvent struct {
	Category    TimelineCategory `json:"category"`
	Title       string           `json:"title"`
	StartDate   string           `json:"start_date,omitempty"`
	EndDate     string           `json:"end_date,omitempty"`
	Description string           `json:"description,omitempty"`
	SourceType  string           `json:"source_type,omitempty"`
	SourceRef   string           `json:"source_ref,omitempty"`
}

// AttributeTimelinesPayloadV1 — payload for type=timelines.
type AttributeTimelinesPayloadV1 struct {
	Events []TimelineEvent `json:"events"`
}

// -----------------------------------------------------------------------------
// Tier 3 — Stubs (post-beta / post-GA)
//
// Real Go types so the schema browser and writers know the eventual shape.
// Expected to be empty / minimally populated until their writers exist.
// -----------------------------------------------------------------------------

// NewsContentType is the kind of news content (article / video / podcast).
type NewsContentType string

const (
	NewsContentTypeArticle NewsContentType = "article"
	NewsContentTypeVideo   NewsContentType = "video"
	NewsContentTypePodcast NewsContentType = "podcast"
)

// NewsMatchConfidence is the per-article match confidence.
type NewsMatchConfidence string

const (
	NewsMatchVerified NewsMatchConfidence = "verified"
	NewsMatchInferred NewsMatchConfidence = "inferred"
)

// NewsArticle is one entry in the in-the-news payload.
type NewsArticle struct {
	Title           string              `json:"title"`
	Source          string              `json:"source"`
	PublishedAt     string              `json:"published_at,omitempty"`
	URL             string              `json:"url"`
	MatchConfidence NewsMatchConfidence `json:"match_confidence"`
	ContentType     NewsContentType     `json:"content_type,omitempty"`
	ImageURL        string              `json:"image_url,omitempty"`
}

// AttributeInTheNewsPayloadV1 — payload for type=in-the-news.
type AttributeInTheNewsPayloadV1 struct {
	Articles []NewsArticle `json:"articles"`
}

// DataLeakSource describes the platform/breach a leak originated from.
type DataLeakSource struct {
	Name   string `json:"name"`
	Domain string `json:"domain,omitempty"`
}

// DataLeakSeverity rates the severity of an exposure.
type DataLeakSeverity string

const (
	DataLeakSeverityLow      DataLeakSeverity = "low"
	DataLeakSeverityMedium   DataLeakSeverity = "medium"
	DataLeakSeverityHigh     DataLeakSeverity = "high"
	DataLeakSeverityCritical DataLeakSeverity = "critical"
)

// DataLeak is one exposure entry. ExposedCategories holds category NAMES only
// (e.g. "email", "password-hash", "ssn-partial") — never the actual leaked values.
type DataLeak struct {
	Source            DataLeakSource   `json:"source"`
	LeakDate          string           `json:"leak_date,omitempty"`
	ExposedCategories []string         `json:"exposed_categories" jsonschema_description:"Category names only (e.g. email, password-hash, ssn-partial); never actual leaked values"`
	Severity          DataLeakSeverity `json:"severity,omitempty"`
}

// AttributeDataLeaksPayloadV1 — payload for type=data-leaks.
//
// Verification-gated: BFF returns null for non-verified callers regardless of
// data availability. Gate is enforced at the BFF layer, not SDS.
type AttributeDataLeaksPayloadV1 struct {
	Leaks []DataLeak `json:"leaks"`
}

// AttributeBehavioralAnalysisPayloadV1 — payload for type=behavioral-analysis.
//
// Skeletal v1. Likely to grow significantly when AE work begins (per-trait
// ratings, evidence per trait, confidence per dimension, etc.); revisit when
// that work starts.
type AttributeBehavioralAnalysisPayloadV1 struct {
	Traits  []string `json:"traits" jsonschema_description:"e.g. 'Avid reader', 'Introvert'"`
	Summary string   `json:"summary,omitempty"`
}

// AttributeSummaryPayloadV1 — payload for type=summary.
//
// kept separate for v1.
type AttributeSummaryPayloadV1 struct {
	Summary string `json:"summary"`
}
