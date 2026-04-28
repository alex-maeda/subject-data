package persistence

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository provides generic JSON-blob CRUD for tables following the standard pattern:
// id (PK), data (JSON TEXT), created_at, updated_at.
type Repository struct {
	db    *sqlx.DB
	table string
}

// NewRepository creates a Repository for the given table.
func NewRepository(db *sqlx.DB, table string) *Repository {
	return &Repository{db: db, table: table}
}

// DB returns the underlying sqlx.DB for direct queries.
func (r *Repository) DB() *sqlx.DB { return r.db }

// Table returns the table name.
func (r *Repository) Table() string { return r.table }

// Upsert inserts or replaces a row with the given id and data object (serialized to JSON).
func (r *Repository) Upsert(id string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("marshaling data: %w", err)
	}
	now := time.Now().UTC()
	query := r.db.Rebind(fmt.Sprintf(
		`INSERT INTO %s (id, data, created_at, updated_at) VALUES (?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET data = excluded.data, updated_at = excluded.updated_at`,
		r.table,
	))
	_, err = r.db.Exec(query, id, string(data), now, now)
	return err
}

// GetByID retrieves a row by primary key, returning the raw JSON data string.
func (r *Repository) GetByID(id string) (string, error) {
	var data string
	query := r.db.Rebind(fmt.Sprintf("SELECT data FROM %s WHERE id = ?", r.table))
	err := r.db.Get(&data, query, id)
	if err != nil {
		return "", fmt.Errorf("getting %s by id %s: %w", r.table, id, err)
	}
	return data, nil
}

// GetAll retrieves all rows, returning raw JSON data strings.
func (r *Repository) GetAll() ([]string, error) {
	var rows []string
	query := fmt.Sprintf("SELECT data FROM %s", r.table)
	err := r.db.Select(&rows, query)
	return rows, err
}

// Delete removes a row by primary key.
func (r *Repository) Delete(id string) error {
	query := r.db.Rebind(fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.table))
	_, err := r.db.Exec(query, id)
	return err
}

// IndexedRepository extends Repository with support for secondary key columns.
// Tables must follow the pattern: id (PK), <key columns>, data (JSON TEXT), created_at, updated_at.
type IndexedRepository struct {
	Repository
	keyColumns []string
}

// NewIndexedRepository creates a repository with the given secondary key columns.
func NewIndexedRepository(db *sqlx.DB, table string, keyColumns ...string) *IndexedRepository {
	return &IndexedRepository{
		Repository: Repository{db: db, table: table},
		keyColumns: keyColumns,
	}
}

// KeyColumns returns the secondary key column names.
func (r *IndexedRepository) KeyColumns() []string { return r.keyColumns }

// UpsertWithKeys inserts or replaces a row with secondary key values.
// keyValues must be in the same order as the key columns passed to NewIndexedRepository.
func (r *IndexedRepository) UpsertWithKeys(id string, keyValues []string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("marshaling data: %w", err)
	}
	now := time.Now().UTC()

	// Build column list: id, key1, key2, ..., data, created_at, updated_at
	cols := append([]string{"id"}, r.keyColumns...)
	cols = append(cols, "data", "created_at", "updated_at")

	placeholders := make([]string, len(cols))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	// Build SET clause for conflict: data, each key, updated_at
	var setClauses []string
	setClauses = append(setClauses, "data = excluded.data")
	for _, col := range r.keyColumns {
		setClauses = append(setClauses, fmt.Sprintf("%s = excluded.%s", col, col))
	}
	setClauses = append(setClauses, "updated_at = excluded.updated_at")

	query := r.db.Rebind(fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT(id) DO UPDATE SET %s",
		r.table,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(setClauses, ", "),
	))

	// Build args: id, key values, data, created_at, updated_at
	args := []interface{}{id}
	for _, v := range keyValues {
		args = append(args, v)
	}
	args = append(args, string(data), now, now)

	_, err = r.db.Exec(query, args...)
	return err
}

// GetByKey retrieves all rows matching a secondary key value.
// keyColumn must be one of the key columns passed to NewIndexedRepository.
func (r *IndexedRepository) GetByKey(keyColumn, value string) ([]string, error) {
	var rows []string
	query := r.db.Rebind(fmt.Sprintf("SELECT data FROM %s WHERE %s = ?", r.table, keyColumn))
	err := r.db.Select(&rows, query, value)
	return rows, err
}

// DeleteByKey removes all rows matching a secondary key value.
func (r *IndexedRepository) DeleteByKey(keyColumn, value string) (int64, error) {
	query := r.db.Rebind(fmt.Sprintf("DELETE FROM %s WHERE %s = ?", r.table, keyColumn))
	result, err := r.db.Exec(query, value)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// --- Convenience constructors and methods ---

// SubjectIndexedRepository is a type alias documenting that an IndexedRepository
// has "subject_id" as its secondary key.
type SubjectIndexedRepository = IndexedRepository

// NewSubjectIndexedRepository creates an IndexedRepository with "subject_id" as the secondary key.
func NewSubjectIndexedRepository(db *sqlx.DB, table string) *IndexedRepository {
	return NewIndexedRepository(db, table, "subject_id")
}

// Upsert is a convenience method for single-key repositories.
func (r *IndexedRepository) Upsert(id, keyValue string, obj interface{}) error {
	return r.UpsertWithKeys(id, []string{keyValue}, obj)
}

// GetBySubjectID is a convenience method for subject-keyed repositories.
func (r *IndexedRepository) GetBySubjectID(subjectID string) ([]string, error) {
	return r.GetByKey("subject_id", subjectID)
}
