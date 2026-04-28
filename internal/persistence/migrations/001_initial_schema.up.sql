CREATE TABLE IF NOT EXISTS subjects (
    id TEXT PRIMARY KEY,
    subject_name TEXT NOT NULL UNIQUE,
    data TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS records (
    id TEXT PRIMARY KEY,
    subject_id TEXT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    data TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS ix_records_subject_id ON records (subject_id);

CREATE TABLE IF NOT EXISTS subject_ratings (
    id TEXT PRIMARY KEY,
    subject_id TEXT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    data TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS ix_subject_ratings_subject_id ON subject_ratings (subject_id);

CREATE TABLE IF NOT EXISTS subject_ce_features (
    id TEXT PRIMARY KEY,
    subject_id TEXT NOT NULL,
    subject_ratings_id TEXT NOT NULL,
    data TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS ix_subject_ce_features_subject_id ON subject_ce_features (subject_id);
CREATE INDEX IF NOT EXISTS ix_subject_ce_features_ratings_id ON subject_ce_features (subject_ratings_id);
