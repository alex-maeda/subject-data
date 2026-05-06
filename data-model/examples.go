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
	}
}
