package datamodel

// RatingBenchmark describes a subject with a particular rating.
type RatingBenchmark struct {
	RatingAsNumber int    `json:"rating_as_number" jsonschema_description:"Rating 1-5"`
	Description    string `json:"description" jsonschema_description:"Description of a subject with a matching rating"`
}

// Evidence is a piece of evidence from a source document.
type Evidence struct {
	Identifiable
	TraceID       string  `json:"trace_id" jsonschema_description:"Identifier of the trace/document this evidence comes from"`
	PageNumber    *string `json:"page_number,omitempty" jsonschema_description:"Page number within the source document, if applicable"`
	Quote         *string `json:"quote,omitempty" jsonschema_description:"Direct quote from the source material"`
	RelevanceNote string  `json:"relevance_note" jsonschema_description:"Explanation of why this evidence is relevant to the rating"`
}

// EvidenceEnriched extends Evidence with dereferenced trace data.
type EvidenceEnriched struct {
	Evidence
	Record *Record `json:"record,omitempty" jsonschema_description:"Dereferenced record data for this evidence"`
}

// CEFeatureDefinition defines a CE feature with benchmarks.
type CEFeatureDefinition struct {
	Identifiable
	Feature     string            `json:"feature" jsonschema_description:"The feature identifier or name"`
	Category    string            `json:"category" jsonschema_description:"Category in which the feature is grouped"`
	Subcategory *string           `json:"subcategory,omitempty" jsonschema_description:"Subcategory within the category"`
	Definition  string            `json:"definition" jsonschema_description:"Description of the feature"`
	Benchmarks  []RatingBenchmark `json:"benchmarks" jsonschema_description:"Description of a subject with different ratings"`
}

// NewCEFeatureDefinition creates a CEFeatureDefinition with initialized slices.
func NewCEFeatureDefinition() CEFeatureDefinition {
	return CEFeatureDefinition{
		Benchmarks: []RatingBenchmark{},
	}
}

// CEFeatureDefinitionEnriched extends CEFeatureDefinition with derived benchmark accessors.
type CEFeatureDefinitionEnriched struct {
	CEFeatureDefinition
	LowBenchmark  *string `json:"low_benchmark,omitempty" jsonschema_description:"Description for a low (1) rating"`
	MedBenchmark  *string `json:"med_benchmark,omitempty" jsonschema_description:"Description for a medium (3) rating"`
	HighBenchmark *string `json:"high_benchmark,omitempty" jsonschema_description:"Description for a high (5) rating"`
}

// EnrichCEFeatureDefinition creates an enriched version with computed benchmark fields.
func EnrichCEFeatureDefinition(def CEFeatureDefinition) CEFeatureDefinitionEnriched {
	e := CEFeatureDefinitionEnriched{CEFeatureDefinition: def}
	for _, b := range def.Benchmarks {
		desc := b.Description
		switch b.RatingAsNumber {
		case 1:
			e.LowBenchmark = &desc
		case 3:
			e.MedBenchmark = &desc
		case 5:
			e.HighBenchmark = &desc
		}
	}
	return e
}

// SubjectCEFeature is a subject-specific CE feature rating with evidence.
type SubjectCEFeature struct {
	Identifiable
	SubjectID           *string    `json:"subject_id,omitempty" jsonschema_description:"Subject identifier"`
	SubjectRatingsID    *string    `json:"subject_ratings_id,omitempty" jsonschema_description:"ID of the SubjectRatings this feature belongs to"`
	Feature             string     `json:"feature" jsonschema_description:"The feature identifier or name"`
	Category            string     `json:"category" jsonschema_description:"Category in which the feature is grouped"`
	Subcategory         *string    `json:"subcategory,omitempty" jsonschema_description:"Subcategory within the category"`
	Rating              string     `json:"rating" jsonschema_description:"Rating as a string: L, L/M, M, M/H, H"`
	RatingAsNumber      *int       `json:"rating_as_number,omitempty" jsonschema_description:"Rating 1-5"`
	RatingReasoning     string     `json:"rating_reasoning" jsonschema_description:"Why the rating was assigned"`
	Evidence            []Evidence `json:"evidence" jsonschema_description:"Trace data sources supporting the reasoning"`
	Confidence          string     `json:"confidence" jsonschema_description:"Confidence as a string: L, M, H"`
	ConfidenceAsNumber  *int       `json:"confidence_as_number,omitempty" jsonschema_description:"Confidence 1/3/5"`
	ConfidenceReasoning *string    `json:"confidence_reasoning,omitempty" jsonschema_description:"Why the confidence level was assigned"`
	Tags                []string   `json:"tags" jsonschema_description:"A list of tags attached to the SubjectCEFeature"`
}

// NewSubjectCEFeature creates a SubjectCEFeature with initialized slices.
func NewSubjectCEFeature() SubjectCEFeature {
	return SubjectCEFeature{
		Evidence: []Evidence{},
		Tags:     []string{},
	}
}

// SubjectRatings is a container for all rating data from a single assessment.
type SubjectRatings struct {
	Identifiable
	SubjectID      *string            `json:"subject_id,omitempty" jsonschema_description:"Subject identifier"`
	RaterID        string             `json:"rater_id" jsonschema_description:"Name or identifier of the rater"`
	Date           *DateTime          `json:"date,omitempty" jsonschema_description:"Date the rating was produced (ISO 8601)"`
	FeatureRatings []SubjectCEFeature `json:"feature_ratings" jsonschema_description:"Subject's feature ratings"`
}

// NewSubjectRatings creates a SubjectRatings with initialized slices.
func NewSubjectRatings() SubjectRatings {
	return SubjectRatings{
		FeatureRatings: []SubjectCEFeature{},
	}
}

// SubjectRatingsEnriched extends SubjectRatings with derived category accessors.
type SubjectRatingsEnriched struct {
	SubjectRatings
	Character              map[string]SubjectCEFeature `json:"character" jsonschema_description:"Character category ratings keyed by feature name"`
	NormalPersonality      map[string]SubjectCEFeature `json:"normal_personality" jsonschema_description:"Normal Personality category ratings keyed by feature name"`
	PersonalityUnderStress map[string]SubjectCEFeature `json:"personality_under_stress" jsonschema_description:"Personality Under Stress category ratings keyed by feature name"`
	Psychopathy            map[string]SubjectCEFeature `json:"psychopathy" jsonschema_description:"Psychopathy category ratings keyed by feature name"`
	Values                 map[string]SubjectCEFeature `json:"values" jsonschema_description:"Values category ratings keyed by feature name"`
	Motivations            map[string]SubjectCEFeature `json:"motivations" jsonschema_description:"Motivations category ratings keyed by feature name"`
	Aptitude               map[string]SubjectCEFeature `json:"aptitude" jsonschema_description:"Aptitude category ratings keyed by feature name"`
	PersonalityDx          map[string]SubjectCEFeature `json:"personality_dx" jsonschema_description:"Personality Dx category ratings keyed by feature name"`
	RiskTaxonomy           map[string]SubjectCEFeature `json:"risk_taxonomy" jsonschema_description:"Risk Taxonomy category ratings keyed by feature name"`
}

// EnrichSubjectRatings creates an enriched version with computed category maps.
func EnrichSubjectRatings(sr SubjectRatings) SubjectRatingsEnriched {
	e := SubjectRatingsEnriched{
		SubjectRatings:         sr,
		Character:              map[string]SubjectCEFeature{},
		NormalPersonality:      map[string]SubjectCEFeature{},
		PersonalityUnderStress: map[string]SubjectCEFeature{},
		Psychopathy:            map[string]SubjectCEFeature{},
		Values:                 map[string]SubjectCEFeature{},
		Motivations:            map[string]SubjectCEFeature{},
		Aptitude:               map[string]SubjectCEFeature{},
		PersonalityDx:          map[string]SubjectCEFeature{},
		RiskTaxonomy:           map[string]SubjectCEFeature{},
	}
	for _, f := range sr.FeatureRatings {
		switch Category(f.Category) {
		case CategoryCharacter:
			e.Character[f.Feature] = f
		case CategoryNormalPersonality:
			e.NormalPersonality[f.Feature] = f
		case CategoryPersonalityUnderStress:
			e.PersonalityUnderStress[f.Feature] = f
		case CategoryPsychopathy:
			e.Psychopathy[f.Feature] = f
		case CategoryValues:
			e.Values[f.Feature] = f
		case CategoryMotivations:
			e.Motivations[f.Feature] = f
		case CategoryAptitude:
			e.Aptitude[f.Feature] = f
		case CategoryPersonalityDx:
			e.PersonalityDx[f.Feature] = f
		case CategoryRiskTaxonomy:
			e.RiskTaxonomy[f.Feature] = f
		}
	}
	return e
}
