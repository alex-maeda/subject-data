package persistence

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// SubjectRepo provides CRUD specific to the subjects table, which has a
// subject_name column in addition to the standard id + data pattern.
type SubjectRepo struct {
	db *sqlx.DB
}

// NewSubjectRepo creates a SubjectRepo.
func NewSubjectRepo(db *sqlx.DB) *SubjectRepo {
	return &SubjectRepo{db: db}
}

// Upsert inserts or replaces a subject row.
func (r *SubjectRepo) Upsert(id, subjectName string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("marshaling data: %w", err)
	}
	now := time.Now().UTC()
	query := r.db.Rebind(
		`INSERT INTO subjects (id, subject_name, data, created_at, updated_at) VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET subject_name = excluded.subject_name, data = excluded.data, updated_at = excluded.updated_at`,
	)
	_, err = r.db.Exec(query, id, subjectName, string(data), now, now)
	return err
}

// GetByID retrieves a subject by primary key.
func (r *SubjectRepo) GetByID(id string) (string, error) {
	var data string
	query := r.db.Rebind("SELECT data FROM subjects WHERE id = ?")
	err := r.db.Get(&data, query, id)
	if err != nil {
		return "", fmt.Errorf("getting subject by id %s: %w", id, err)
	}
	return data, nil
}

// GetAll retrieves all subjects.
func (r *SubjectRepo) GetAll() ([]string, error) {
	var rows []string
	err := r.db.Select(&rows, "SELECT data FROM subjects")
	return rows, err
}

// Delete removes a subject by primary key.
func (r *SubjectRepo) Delete(id string) error {
	query := r.db.Rebind("DELETE FROM subjects WHERE id = ?")
	_, err := r.db.Exec(query, id)
	return err
}
