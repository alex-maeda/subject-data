package datamodel

// Category in which a feature is grouped.
type Category string

// Category constants.
const (
	CategoryNormalPersonality      Category = "Normal Personality"
	CategoryPersonalityUnderStress Category = "Personality Under Stress"
	CategoryCharacter              Category = "Character"
	CategoryPsychopathy            Category = "Psychopathy"
	CategoryValues                 Category = "Values"
	CategoryMotivations            Category = "Motivations"
	CategoryAptitude               Category = "Aptitude"
	CategoryRiskTaxonomy           Category = "Risk Taxonomy"
	CategoryPersonalityDx          Category = "Personality Dx"
)

// Subcategory within a category.
type Subcategory string

// Subcategory constants.
const (
	SubcategoryBigFive          Subcategory = "Big Five"
	SubcategoryHEXACO           Subcategory = "HEXACO"
	SubcategoryMainRiskCategory Subcategory = "Main Risk Category"
)

// RatingStr is a rating as a string.
type RatingStr string

// RatingStr constants.
const (
	RatingLow        RatingStr = "Low"
	RatingLowMedium  RatingStr = "Low to Moderate"
	RatingMedium     RatingStr = "Moderate"
	RatingMediumHigh RatingStr = "Moderate to High"
	RatingHigh       RatingStr = "High"
)

// ConfidenceStr is a confidence as a string.
type ConfidenceStr string

// ConfidenceStr constants.
const (
	ConfidenceLow    ConfidenceStr = "Low"
	ConfidenceMedium ConfidenceStr = "Medium"
	ConfidenceHigh   ConfidenceStr = "High"
)

// Platform is the source platform for trace data.
type Platform string

// Platform constants.
const (
	PlatformBluesky               Platform = "bluesky"
	PlatformCrunchbase            Platform = "crunchbase"
	PlatformFacebook              Platform = "facebook"
	PlatformFlightradar24         Platform = "flightradar24"
	PlatformGlassdoor             Platform = "glassdoor"
	PlatformGoogle                Platform = "google"
	PlatformInstagram             Platform = "instagram"
	PlatformLinkedin              Platform = "linkedin"
	PlatformReddit                Platform = "reddit"
	PlatformThreads               Platform = "threads"
	PlatformTiktok                Platform = "tiktok"
	PlatformTwitter               Platform = "twitter"
	PlatformSkopenow              Platform = "skopenow"
	PlatformTransunion            Platform = "transunion"
	PlatformBespokeCriminalRecord Platform = "bespoke_criminal_record"
)

// ContentType classifies the content type, shared across platforms.
type ContentType string

// ContentType constants.
const (
	ContentTypePost                 ContentType = "post"
	ContentTypeProfile              ContentType = "profile"
	ContentTypeComment              ContentType = "comment"
	ContentTypeReaction             ContentType = "reaction"
	ContentTypeDocument             ContentType = "document"
	ContentTypeBankruptcyRecords    ContentType = "bankruptcy_records"
	ContentTypeCivilCourtRecords    ContentType = "civil_court_records"
	ContentTypeCriminalCourtRecords ContentType = "criminal_court_records"
	ContentTypeFamilyCourtRecords   ContentType = "family_court_records"
	ContentTypeMortgageDeedRecords  ContentType = "mortgage_deed_records"
)

// ContextType describes how an image or video was used in the trace.
type ContextType string

// ContextType constants.
const (
	ContextTypeDefault  ContextType = "default"
	ContextTypeOriginal ContextType = "original"
	ContextTypeReaction ContextType = "reaction"
)
