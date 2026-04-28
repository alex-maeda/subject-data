// Package service provides business logic for subjects, records, and ratings.
package service

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	datamodel "github.com/sovraai/subject-data/data-model"
	"github.com/sovraai/subject-data/internal/persistence"
)

// RecordService provides CRUD operations for records (formerly TraceData).
type RecordService struct {
	logger *slog.Logger
	repo   *persistence.SubjectIndexedRepository
}

// NewRecordService creates a RecordService backed by the given database.
func NewRecordService(logger *slog.Logger, db *sqlx.DB) *RecordService {
	return &RecordService{
		logger: logger,
		repo:   persistence.NewSubjectIndexedRepository(db, "records"),
	}
}

// Create persists a new record and returns its ID.
func (s *RecordService) Create(obj *datamodel.Record) (string, error) {
	id := ""
	if obj.ID != nil {
		id = *obj.ID
	}
	if id == "" {
		id = uuid.New().String()
		obj.ID = &id
	}
	subjectID := ""
	if obj.SubjectID != nil {
		subjectID = *obj.SubjectID
	}
	if err := s.repo.Upsert(id, subjectID, obj); err != nil {
		return "", fmt.Errorf("creating record: %w", err)
	}
	return id, nil
}

// Update replaces an existing record by ID.
func (s *RecordService) Update(recordID string, obj *datamodel.Record) error {
	if _, err := s.repo.GetByID(recordID); err != nil {
		return fmt.Errorf("record not found: %s", recordID)
	}
	obj.ID = &recordID
	subjectID := ""
	if obj.SubjectID != nil {
		subjectID = *obj.SubjectID
	}
	return s.repo.Upsert(recordID, subjectID, obj)
}

// Get retrieves a record by ID.
func (s *RecordService) Get(recordID string) (*datamodel.Record, error) {
	data, err := s.repo.GetByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("record not found: %s", recordID)
	}
	var obj datamodel.Record
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return nil, fmt.Errorf("unmarshaling record: %w", err)
	}
	obj.ID = &recordID
	return &obj, nil
}

// Delete removes a record by ID.
func (s *RecordService) Delete(recordID string) error {
	return s.repo.Delete(recordID)
}

// ListAll returns all records.
func (s *RecordService) ListAll() ([]datamodel.Record, error) {
	var rows []persistence.SubjectKeyedRow
	err := s.repo.DB().Select(&rows, "SELECT id, subject_id, data FROM records")
	if err != nil {
		return nil, err
	}
	return unmarshalRecordRows(rows)
}

// ListBySubject returns all records for a given subject.
func (s *RecordService) ListBySubject(subjectID string) ([]datamodel.Record, error) {
	var rows []persistence.SubjectKeyedRow
	query := s.repo.DB().Rebind("SELECT id, subject_id, data FROM records WHERE subject_id = ?")
	err := s.repo.DB().Select(&rows, query, subjectID)
	if err != nil {
		return nil, err
	}
	return unmarshalRecordRows(rows)
}

func unmarshalRecordRows(rows []persistence.SubjectKeyedRow) ([]datamodel.Record, error) {
	result := make([]datamodel.Record, 0, len(rows))
	for _, row := range rows {
		var obj datamodel.Record
		if err := json.Unmarshal([]byte(row.Data), &obj); err != nil {
			return nil, fmt.Errorf("unmarshaling record row %s: %w", row.ID, err)
		}
		id := row.ID
		obj.ID = &id
		result = append(result, obj)
	}
	return result, nil
}
