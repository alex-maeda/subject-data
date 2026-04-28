package persistence

import "time"

// StandardRow is the common pattern for most tables: id, data JSON blob, timestamps.
type StandardRow struct {
	ID        string    `db:"id"`
	Data      string    `db:"data"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// SubjectRow stores subject-level data as a JSON blob.
type SubjectRow struct {
	ID          string    `db:"id"`
	SubjectName string    `db:"subject_name"`
	Data        string    `db:"data"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// SubjectKeyedRow is the pattern for tables with id + subject_id + data.
type SubjectKeyedRow struct {
	ID        string    `db:"id"`
	SubjectID string    `db:"subject_id"`
	Data      string    `db:"data"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// SubjectCEFeatureRow adds subject_ratings_id to the subject-keyed pattern.
type SubjectCEFeatureRow struct {
	ID               string    `db:"id"`
	SubjectID        string    `db:"subject_id"`
	SubjectRatingsID *string   `db:"subject_ratings_id"`
	Data             string    `db:"data"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
