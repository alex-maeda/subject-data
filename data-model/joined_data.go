package datamodel

// JoinedSubjectData combines ratings and input data for a single subject.
type JoinedSubjectData struct {
	ID          *string          `json:"id,omitempty" jsonschema_description:"Opaque subject identifier"`
	SubjectID   *string          `json:"subject_id,omitempty" jsonschema_description:"Subject identifier"`
	SubjectName *string          `json:"subject_name,omitempty" jsonschema_description:"Subject's name"`
	TraceData   []Record         `json:"trace_data" jsonschema_description:"Data from processed reports and scraped data"`
	Ratings     []SubjectRatings `json:"ratings" jsonschema_description:"List of assessment ratings, one per rater"`
}

// NewJoinedSubjectData creates a JoinedSubjectData with initialized slices.
func NewJoinedSubjectData() JoinedSubjectData {
	return JoinedSubjectData{
		TraceData: []Record{},
		Ratings:   []SubjectRatings{},
	}
}

// JoinedDataset is a complete joined dataset containing multiple subjects.
type JoinedDataset struct {
	Examples []JoinedSubjectData `json:"examples" jsonschema_description:"List of subject data entries"`
}

// NewJoinedDataset creates a JoinedDataset with initialized slices.
func NewJoinedDataset() JoinedDataset {
	return JoinedDataset{
		Examples: []JoinedSubjectData{},
	}
}

// JoinedDatasetEnriched extends JoinedDataset with derived subject/rater accessors.
type JoinedDatasetEnriched struct {
	JoinedDataset
	Subjects []string `json:"subjects" jsonschema_description:"List of unique subject IDs in the dataset"`
	Raters   []string `json:"raters" jsonschema_description:"List of unique rater IDs in the dataset"`
}

// EnrichJoinedDataset creates an enriched version with computed subject and rater lists.
func EnrichJoinedDataset(ds JoinedDataset) JoinedDatasetEnriched {
	subjectSet := map[string]struct{}{}
	raterSet := map[string]struct{}{}
	for _, e := range ds.Examples {
		if e.SubjectID != nil && *e.SubjectID != "" {
			subjectSet[*e.SubjectID] = struct{}{}
		}
		for _, r := range e.Ratings {
			if r.RaterID != "" {
				raterSet[r.RaterID] = struct{}{}
			}
		}
	}
	subjects := make([]string, 0, len(subjectSet))
	for s := range subjectSet {
		subjects = append(subjects, s)
	}
	raters := make([]string, 0, len(raterSet))
	for r := range raterSet {
		raters = append(raters, r)
	}
	return JoinedDatasetEnriched{
		JoinedDataset: ds,
		Subjects:      subjects,
		Raters:        raters,
	}
}

// ExtractTraitScores extracts numeric trait scores and reasoning from an assessment.
// Returns (traitScores, traitReasoning) keyed by "category:trait_name".
func ExtractTraitScores(assessment SubjectRatings) (map[string]*int, map[string]string) {
	scores := map[string]*int{}
	reasoning := map[string]string{}

	for _, facet := range assessment.FeatureRatings {
		fullName := facet.Category + ":" + facet.Feature
		scores[fullName] = facet.RatingAsNumber

		if facet.RatingReasoning != "" {
			reasoning[fullName] = facet.RatingReasoning
		}
	}

	return scores, reasoning
}
