package service

import (
	"encoding/json"
	"fmt"
	"time"

	datamodel "github.com/sovraai/subject-data/data-model"
	"github.com/sovraai/subject-data/internal/persistence"
)

// RatingStorageService provides persistence for subject ratings and CE features.
type RatingStorageService struct {
	ratings  *persistence.SubjectIndexedRepository
	features *persistence.SubjectIndexedRepository
}

// NewRatingStorageService creates a RatingStorageService.
func NewRatingStorageService(
	ratings *persistence.SubjectIndexedRepository,
	features *persistence.SubjectIndexedRepository,
) *RatingStorageService {
	return &RatingStorageService{
		ratings:  ratings,
		features: features,
	}
}

// SaveRatings persists a SubjectRatings and its feature ratings.
func (s *RatingStorageService) SaveRatings(subjectID string, sr *datamodel.SubjectRatings) error {
	if sr.ID == nil {
		return fmt.Errorf("ratings ID is required")
	}
	ratingsID := *sr.ID

	// Save the ratings container (with empty feature_ratings in the row)
	container := *sr
	container.FeatureRatings = []datamodel.SubjectCEFeature{}
	if err := s.ratings.Upsert(ratingsID, subjectID, container); err != nil {
		return fmt.Errorf("saving ratings: %w", err)
	}

	// Save each feature rating with subject_ratings_id column set.
	for _, f := range sr.FeatureRatings {
		if f.ID == nil {
			continue
		}
		f.SubjectRatingsID = &ratingsID
		data, err := json.Marshal(f)
		if err != nil {
			return fmt.Errorf("marshaling feature: %w", err)
		}
		now := time.Now().UTC()
		query := s.features.DB().Rebind(
			`INSERT INTO subject_ce_features (id, subject_id, subject_ratings_id, data, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)
			 ON CONFLICT(id) DO UPDATE SET data = excluded.data, subject_id = excluded.subject_id, subject_ratings_id = excluded.subject_ratings_id, updated_at = excluded.updated_at`,
		)
		_, err = s.features.DB().Exec(query, *f.ID, subjectID, ratingsID, string(data), now, now)
		if err != nil {
			return fmt.Errorf("saving feature: %w", err)
		}
	}
	return nil
}

// GetRatingsBySubject retrieves all ratings for a subject, with features loaded.
func (s *RatingStorageService) GetRatingsBySubject(subjectID string) ([]datamodel.SubjectRatings, error) {
	rows, err := s.ratings.GetBySubjectID(subjectID)
	if err != nil {
		return nil, fmt.Errorf("getting ratings: %w", err)
	}

	result := make([]datamodel.SubjectRatings, 0, len(rows))
	for _, row := range rows {
		var sr datamodel.SubjectRatings
		if err := json.Unmarshal([]byte(row), &sr); err != nil {
			return nil, fmt.Errorf("unmarshaling ratings: %w", err)
		}

		if sr.ID != nil {
			features, err := s.loadFeaturesByRatingsID(*sr.ID)
			if err != nil {
				return nil, fmt.Errorf("loading features for ratings %s: %w", *sr.ID, err)
			}
			sr.FeatureRatings = features
		}

		result = append(result, sr)
	}
	return result, nil
}

// loadFeaturesByRatingsID queries subject_ce_features for a given subject_ratings_id.
func (s *RatingStorageService) loadFeaturesByRatingsID(ratingsID string) ([]datamodel.SubjectCEFeature, error) {
	var featureRows []persistence.SubjectCEFeatureRow
	query := s.features.DB().Rebind(
		"SELECT id, subject_id, subject_ratings_id, data, created_at, updated_at FROM subject_ce_features WHERE subject_ratings_id = ?",
	)
	err := s.features.DB().Select(&featureRows, query, ratingsID)
	if err != nil {
		return nil, err
	}

	features := make([]datamodel.SubjectCEFeature, 0, len(featureRows))
	for _, row := range featureRows {
		var f datamodel.SubjectCEFeature
		if err := json.Unmarshal([]byte(row.Data), &f); err != nil {
			return nil, fmt.Errorf("unmarshaling feature %s: %w", row.ID, err)
		}
		id := row.ID
		f.ID = &id
		features = append(features, f)
	}
	return features, nil
}
