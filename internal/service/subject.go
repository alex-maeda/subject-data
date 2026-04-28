package service

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	datamodel "github.com/sovraai/subject-data/data-model"
	"github.com/sovraai/subject-data/internal/persistence"
)

// SubjectService provides CRUD operations for subjects with their associated
// records, ratings, and CE features.
type SubjectService struct {
	logger   *slog.Logger
	subjects *persistence.SubjectRepo
	records  *persistence.SubjectIndexedRepository
	ratings  *persistence.SubjectIndexedRepository
	features *persistence.SubjectIndexedRepository
}

// NewSubjectService creates a SubjectService backed by the given repositories.
func NewSubjectService(
	logger *slog.Logger,
	subjects *persistence.SubjectRepo,
	records *persistence.SubjectIndexedRepository,
	ratings *persistence.SubjectIndexedRepository,
	features *persistence.SubjectIndexedRepository,
) *SubjectService {
	return &SubjectService{
		logger:   logger,
		subjects: subjects,
		records:  records,
		ratings:  ratings,
		features: features,
	}
}

// GetSubject retrieves a subject by ID with all associated data.
func (s *SubjectService) GetSubject(id string) (*datamodel.JoinedSubjectData, error) {
	data, err := s.subjects.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("getting subject: %w", err)
	}
	var subject datamodel.JoinedSubjectData
	if err := json.Unmarshal([]byte(data), &subject); err != nil {
		return nil, fmt.Errorf("unmarshaling subject: %w", err)
	}
	subject.SubjectID = &id

	// Load records
	recordRows, err := s.records.GetBySubjectID(id)
	if err != nil {
		return nil, fmt.Errorf("getting records: %w", err)
	}
	subject.TraceData = make([]datamodel.Record, 0, len(recordRows))
	for _, row := range recordRows {
		var td datamodel.Record
		if err := json.Unmarshal([]byte(row), &td); err != nil {
			return nil, fmt.Errorf("unmarshaling record: %w", err)
		}
		subject.TraceData = append(subject.TraceData, td)
	}

	// Load ratings with nested feature_ratings
	ratingSvc := NewRatingStorageService(s.ratings, s.features)
	ratings, err := ratingSvc.GetRatingsBySubject(id)
	if err != nil {
		return nil, fmt.Errorf("getting ratings: %w", err)
	}
	subject.Ratings = ratings

	return &subject, nil
}

// ListSubjects returns all subjects (without records or ratings).
func (s *SubjectService) ListSubjects() ([]datamodel.JoinedSubjectData, error) {
	dataRows, err := s.subjects.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]datamodel.JoinedSubjectData, 0, len(dataRows))
	for _, row := range dataRows {
		var subject datamodel.JoinedSubjectData
		if err := json.Unmarshal([]byte(row), &subject); err != nil {
			return nil, fmt.Errorf("unmarshaling subject: %w", err)
		}
		result = append(result, subject)
	}
	return result, nil
}

// CreateSubject creates a new subject with the given name. Returns the new subject's ID.
func (s *SubjectService) CreateSubject(subjectName string) (string, error) {
	id := uuid.New().String()
	subject := &datamodel.JoinedSubjectData{
		SubjectName: &subjectName,
		TraceData:   []datamodel.Record{},
		Ratings:     []datamodel.SubjectRatings{},
	}
	subject.ID = &id
	if err := s.SaveSubject(subject); err != nil {
		return "", fmt.Errorf("creating subject: %w", err)
	}
	return id, nil
}

// DeleteSubject removes a subject by ID. Cascading foreign keys handle associated data.
func (s *SubjectService) DeleteSubject(id string) error {
	if _, err := s.subjects.GetByID(id); err != nil {
		return fmt.Errorf("subject not found: %s", id)
	}
	if err := s.subjects.Delete(id); err != nil {
		return fmt.Errorf("deleting subject: %w", err)
	}
	return nil
}

// SaveSubject persists a subject and its associated records and ratings.
func (s *SubjectService) SaveSubject(subject *datamodel.JoinedSubjectData) error {
	if subject.ID == nil {
		return fmt.Errorf("subject ID is required")
	}
	id := *subject.ID
	subjectID := id
	if subject.SubjectID != nil {
		subjectID = *subject.SubjectID
	}

	// Save subject (without records and ratings in the blob)
	stripped := datamodel.JoinedSubjectData{
		ID:          subject.ID,
		SubjectID:   subject.SubjectID,
		SubjectName: subject.SubjectName,
		TraceData:   []datamodel.Record{},
		Ratings:     []datamodel.SubjectRatings{},
	}
	subjectName := ""
	if subject.SubjectName != nil {
		subjectName = *subject.SubjectName
	}
	if err := s.subjects.Upsert(id, subjectName, stripped); err != nil {
		return fmt.Errorf("saving subject: %w", err)
	}

	// Save records
	for _, td := range subject.TraceData {
		if td.ID == nil {
			continue
		}
		if err := s.records.Upsert(*td.ID, subjectID, td); err != nil {
			return fmt.Errorf("saving record: %w", err)
		}
	}

	// Save ratings
	for _, sr := range subject.Ratings {
		if sr.ID == nil {
			continue
		}
		if err := s.ratings.Upsert(*sr.ID, subjectID, sr); err != nil {
			return fmt.Errorf("saving rating: %w", err)
		}
	}

	return nil
}
