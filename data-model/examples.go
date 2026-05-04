package datamodel

import "time"

// ptr returns a pointer to the given value.
func ptr[T any](v T) *T { return &v }

// --- Reusable helper variables ---

var exampleEvidence = Evidence{
	Identifiable:  Identifiable{ID: ptr("A1")},
	TraceID:       "instagram-jdoe",
	Quote:         ptr("Just completed my third marathon this year!"),
	RelevanceNote: "Demonstrates high achievement motivation and discipline.",
}

var exampleCEFeature = SubjectCEFeature{
	Identifiable:        Identifiable{ID: ptr("feat-001")},
	SubjectID:           ptr("subj_123456"),
	Feature:             "Courage",
	Category:            "Character",
	Rating:              "M/H",
	RatingAsNumber:      ptr(4),
	RatingReasoning:     "Continually perseveres towards goals, even in the face of failure. Two failed congressional bids, now pursuing state office.",
	Evidence:            []Evidence{exampleEvidence},
	Confidence:          "M",
	ConfidenceAsNumber:  ptr(3),
	ConfidenceReasoning: ptr("Very public well documented failures and continued pursuits."),
	Tags:                []string{},
}

var exampleRatingBenchmarkLow = RatingBenchmark{
	RatingAsNumber: 1,
	Description:    "Rarely shows concern about future events or outcomes.",
}

var exampleRatingBenchmarkMed = RatingBenchmark{
	RatingAsNumber: 3,
	Description:    "Sometimes worries about important Factors, but generally manages anxiety well.",
}

var exampleRatingBenchmarkHigh = RatingBenchmark{
	RatingAsNumber: 5,
	Description:    "Frequently worries about multiple Factors and has difficulty controlling anxious thoughts.",
}

var exampleCEFeatureDefinition = CEFeatureDefinition{
	Identifiable: Identifiable{ID: ptr("12345")},
	Feature:      "Courage",
	Category:     "Character",
	Definition:   "Sustained perseverance under prolonged adversity, requiring willpower, grit, and determination despite repeated setbacks.",
	Benchmarks: []RatingBenchmark{
		exampleRatingBenchmarkLow,
		exampleRatingBenchmarkMed,
		exampleRatingBenchmarkHigh,
	},
}

var exampleMediaAnalysis = MediaAnalysis{
	Description:       ptr("Instagram post featuring an individual outdoors holding a sign"),
	Subjects:          ptr("A person wearing a hoodie and jacket, face partially obscured by a cap"),
	Setting:           ptr("Outdoors in a parking lot or public area, overcast sky"),
	ApparentActivity:  ptr("Holding a handwritten sign for a photo op"),
	PresentationStyle: ptr("Casual, candid street-style photo"),
	MoodTone:          ptr("Lighthearted, humorous"),
	TextInMedia:       ptr("Handwritten sign reads: 'Need Money 4 Jordans'"),
}

var exampleRecordMetadata = RecordMetadata{
	Website:     ptr("instagram.com"),
	URL:         ptr("https://www.instagram.com/p/DTWBVT_EilZ"),
	CaptureDate: &DateTime{Time: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)},
	Tags:        []string{},
	Identifiers: []string{},
}

var exampleDataSource = DataSource{
	Platform:    "instagram",
	ContentType: "post",
}

var exampleImageEntry = ImageEntry{
	URL:      "https://hosted-copy-of-data/image_abcdef.jpg",
	Context:  "default",
	Analysis: &exampleMediaAnalysis,
}

var exampleVideoEntry = VideoEntry{
	URL:      "https://hosted-copy-of-data/video_123456.mp4",
	Context:  "default",
	Analysis: &exampleMediaAnalysis,
}

var examplePlatformData = PlatformData{
	Normalized: map[string]interface{}{
		"text":           "Just completed my third marathon this year! #running #goals",
		"reaction_count": 42,
	},
	Raw: map[string]interface{}{
		"page_type":    "post",
		"username":     "jdoe",
		"caption_text": "Just completed my third marathon this year! #running #goals",
		"like_count":   42,
	},
}

var exampleRecord = Record{
	Identifiable: Identifiable{ID: ptr("12345")},
	SubjectID:    ptr("subj_123456"),
	Metadata:     exampleRecordMetadata,
	DataSource:   exampleDataSource,
	PlatformData: examplePlatformData,
	Images:       []ImageEntry{exampleImageEntry},
	Videos:       []VideoEntry{},
}

var exampleSubjectRatings = SubjectRatings{
	Identifiable:   Identifiable{ID: ptr("12345")},
	SubjectID:      ptr("subj_123456"),
	RaterID:        "llm_rater",
	Date:           &DateTime{Time: time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)},
	FeatureRatings: []SubjectCEFeature{exampleCEFeature},
}

// Examples returns a map from type name to an illustrative example instance.
// These are used by the schema browser to show sample values in Swagger UI,
// and can serve as test stubs.
func Examples() map[string]any {
	return map[string]any{
		"RatingBenchmark": exampleRatingBenchmarkHigh,

		"Evidence": exampleEvidence,

		"EvidenceEnriched": EvidenceEnriched{
			Evidence: exampleEvidence,
			Record:   nil,
		},

		"CEFeatureDefinition": exampleCEFeatureDefinition,

		"CEFeatureDefinitionEnriched": CEFeatureDefinitionEnriched{
			CEFeatureDefinition: exampleCEFeatureDefinition,
			LowBenchmark:        ptr("Rarely shows concern about future events or outcomes."),
			MedBenchmark:        ptr("Sometimes worries about important Factors, but generally manages anxiety well."),
			HighBenchmark:       ptr("Frequently worries about multiple Factors and has difficulty controlling anxious thoughts."),
		},

		"SubjectCEFeature": exampleCEFeature,

		"SubjectRatings": exampleSubjectRatings,

		"SubjectRatingsEnriched": SubjectRatingsEnriched{
			SubjectRatings: exampleSubjectRatings,
			Character: map[string]SubjectCEFeature{
				"Courage": exampleCEFeature,
			},
			NormalPersonality:      map[string]SubjectCEFeature{},
			PersonalityUnderStress: map[string]SubjectCEFeature{},
			Psychopathy:            map[string]SubjectCEFeature{},
			Values:                 map[string]SubjectCEFeature{},
			Motivations:            map[string]SubjectCEFeature{},
			Aptitude:               map[string]SubjectCEFeature{},
			PersonalityDx:          map[string]SubjectCEFeature{},
			RiskTaxonomy:           map[string]SubjectCEFeature{},
		},

		"Subject": Subject{
			Identifiable: Identifiable{ID: ptr("12345")},
			SubjectName:  "Smith, John",
		},

		"RecordMetadata": exampleRecordMetadata,
		"DataSource":     exampleDataSource,
		"MediaAnalysis":  exampleMediaAnalysis,
		"ImageEntry":     exampleImageEntry,
		"VideoEntry":     exampleVideoEntry,
		"PlatformData":   examplePlatformData,

		"Record": exampleRecord,

		"JoinedSubjectData": JoinedSubjectData{
			ID:          ptr("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
			SubjectID:   ptr("subj_123456"),
			SubjectName: ptr("Doe, John"),
			Records:     []Record{exampleRecord},
			Ratings:     []SubjectRatings{exampleSubjectRatings},
		},

		"JoinedDatasetEnriched": JoinedDatasetEnriched{
			JoinedDataset: JoinedDataset{
				Examples: []JoinedSubjectData{
					{
						ID:          ptr("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
						SubjectID:   ptr("subj_123456"),
						SubjectName: ptr("Doe, John"),
						Records:     []Record{exampleRecord},
						Ratings:     []SubjectRatings{exampleSubjectRatings},
					},
				},
			},
			Subjects: []string{"subj_123456"},
			Raters:   []string{"llm_rater"},
		},

		// ----- Attribute envelope & provenance -----

		"Attribute": Attribute{
			Identifiable:    Identifiable{ID: ptr("attr_001")},
			SubjectID:       "subj_123456",
			Type:            AttributeTypeName,
			SchemaVersion:   "v1",
			Payload:         []byte(`{"first":"John","last":"Smith"}`),
			Confidence:      0.92,
			EvidenceSummary: "Derived from PDL record and voter-registration match.",
			EvidenceRefs: []AttributeEvidenceRef{
				{RecordID: "rec_abc", SourceField: "full_name", NormalizedAttribute: "first_name"},
			},
		},

		"AttributeEvidenceRef": AttributeEvidenceRef{
			RecordID:            "rec_abc",
			SourceField:         "full_name",
			NormalizedAttribute: "first_name",
		},

		// ----- Tier 1 — Header payloads -----

		"AttributeNamePayloadV1": AttributeNamePayloadV1{
			First:       "John",
			Last:        "Smith",
			DisplayName: "John Smith",
		},

		"AttributeAgePayloadV1": AttributeAgePayloadV1{
			YearOfBirth: 1990,
			BirthDate:   "1990-06-15",
		},

		"AttributeGenderPayloadV1": AttributeGenderPayloadV1{
			Value: GenderMale,
		},

		"AttributePhonePayloadV1": AttributePhonePayloadV1{
			Obfuscated: "+1 713 ***-**89",
			Country:    "US",
		},

		"AttributeEmailPayloadV1": AttributeEmailPayloadV1{
			Obfuscated: "jo********@gmail.com",
		},

		"AttributeAvatarPayloadV1": AttributeAvatarPayloadV1{
			URL: "https://hosted-copy-of-data/avatar_abc123.jpg",
		},

		// ----- Tier 2 — SO-written card payloads -----

		"AttributeRelationshipPayloadV1": AttributeRelationshipPayloadV1{
			Status: RelationshipStatusMarried,
		},

		"SocialAccount": SocialAccount{
			Platform:        SocialPlatformLinkedIn,
			Handle:          "johnsmith",
			ProfileURL:      "https://linkedin.com/in/johnsmith",
			MatchConfidence: SocialMatchVerified,
		},

		"AttributeSocialMediaFootprintPayloadV1": AttributeSocialMediaFootprintPayloadV1{
			Accounts: []SocialAccount{
				{Platform: SocialPlatformLinkedIn, Handle: "johnsmith", ProfileURL: "https://linkedin.com/in/johnsmith", MatchConfidence: SocialMatchVerified},
				{Platform: SocialPlatformTwitter, Handle: "jsmith42", ProfileURL: "https://twitter.com/jsmith42", MatchConfidence: SocialMatchInferred},
			},
		},

		"GeoCoordinates": GeoCoordinates{Lat: 29.7604, Lng: -95.3698},

		"CityLocation": CityLocation{
			City:        "Houston",
			State:       "TX",
			Coordinates: &GeoCoordinates{Lat: 29.7604, Lng: -95.3698},
		},

		"LocationFrequency": LocationFrequency{
			City: "Houston", State: "TX", Count: 12,
			Coordinates: &GeoCoordinates{Lat: 29.7604, Lng: -95.3698},
		},

		"GeographicDataPoint": GeographicDataPoint{
			Location:   GeographicLocation{City: "Houston", State: "TX"},
			Date:       "2025-03-15",
			SourceType: "voter_registration",
		},

		"AttributeGeographicFootprintPayloadV1": AttributeGeographicFootprintPayloadV1{
			LatestLocation: &CityLocation{City: "Houston", State: "TX"},
			LocationFrequency: []LocationFrequency{
				{City: "Houston", State: "TX", Count: 12},
				{City: "Austin", State: "TX", Count: 3},
			},
			DataPoints: []GeographicDataPoint{
				{Location: GeographicLocation{City: "Houston", State: "TX"}, Date: "2025-03-15", SourceType: "voter_registration"},
			},
		},

		// ----- Tier 2 — DE-written card payloads -----

		"PublicRecordParty": PublicRecordParty{Role: "defendant"},

		"PublicRecord": PublicRecord{
			Title:          "Smith v. State of Texas",
			Date:           "2024-08-12",
			Location:       "Harris County, TX",
			Description:    "Traffic violation",
			Role:           "defendant",
			Classification: "dismissed",
			CaseType:       "traffic",
			Disposition:    "case dismissed",
		},

		"PublicRecordGroup": PublicRecordGroup{
			Type: PublicRecordTypeCivil,
			Records: []PublicRecord{
				{Title: "Smith v. State of Texas", Date: "2024-08-12", Location: "Harris County, TX", CaseType: "traffic", Disposition: "case dismissed"},
			},
		},

		"AttributePublicRecordsPayloadV1": AttributePublicRecordsPayloadV1{
			Groups: []PublicRecordGroup{
				{Type: PublicRecordTypeCivil, Records: []PublicRecord{
					{Title: "Smith v. State of Texas", Date: "2024-08-12", Location: "Harris County, TX"},
				}},
			},
		},

		"TimelineEvent": TimelineEvent{
			Category:  TimelineCategoryLocation,
			Title:     "Moved to Houston, TX",
			StartDate: "2022-06-01",
		},

		"AttributeTimelinesPayloadV1": AttributeTimelinesPayloadV1{
			Events: []TimelineEvent{
				{Category: TimelineCategoryLocation, Title: "Moved to Houston, TX", StartDate: "2022-06-01"},
				{Category: TimelineCategoryProfessional, Title: "Started at Acme Corp", StartDate: "2023-01-15"},
			},
		},

		// ----- Tier 3 — Stub card payloads -----

		"NewsArticle": NewsArticle{
			Title:           "Local engineer wins innovation award",
			Source:          "Houston Chronicle",
			PublishedAt:     "2025-11-20",
			URL:             "https://example.com/article/123",
			MatchConfidence: NewsMatchVerified,
			ContentType:     NewsContentTypeArticle,
		},

		"AttributeInTheNewsPayloadV1": AttributeInTheNewsPayloadV1{
			Articles: []NewsArticle{
				{Title: "Local engineer wins innovation award", Source: "Houston Chronicle", PublishedAt: "2025-11-20", URL: "https://example.com/article/123", MatchConfidence: NewsMatchVerified, ContentType: NewsContentTypeArticle},
			},
		},

		"DataLeakSource": DataLeakSource{Name: "ExampleBreach2024", Domain: "example.com"},

		"DataLeak": DataLeak{
			Source:            DataLeakSource{Name: "ExampleBreach2024", Domain: "example.com"},
			LeakDate:          "2024-02-14",
			ExposedCategories: []string{"email", "password-hash"},
			Severity:          DataLeakSeverityMedium,
		},

		"AttributeDataLeaksPayloadV1": AttributeDataLeaksPayloadV1{
			Leaks: []DataLeak{
				{Source: DataLeakSource{Name: "ExampleBreach2024", Domain: "example.com"}, LeakDate: "2024-02-14", ExposedCategories: []string{"email", "password-hash"}, Severity: DataLeakSeverityMedium},
			},
		},

		"AttributeBehavioralAnalysisPayloadV1": AttributeBehavioralAnalysisPayloadV1{
			Traits:  []string{"Achievement-oriented", "Detail-conscious", "Resilient under pressure"},
			Summary: "Subject demonstrates strong goal-directed behavior with high persistence.",
		},

		"AttributeSummaryPayloadV1": AttributeSummaryPayloadV1{
			Summary: "John Smith is a Houston-based professional with a clean public record and active social media presence.",
		},
	}
}
