-- Migration 003: Attributes table.
--
-- Attributes are typed, schema'd values attached to Subjects. Each AttributeType
-- (name, phone, public-records, ...) has its own payload shape stored in `data`
-- as JSON; the envelope's identity + key + audit fields live as columns to
-- support PK and secondary-index lookups, in line with the project's NoSQL
-- portability constraint.
--
-- Constraints:
--   - UNIQUE (subject_id, type) enforces "one Attribute per (subject, type) pair".
--     Multi-writer reconciliation is implicit via UPSERT (latest write wins).
--   - ON DELETE CASCADE on subject_id: Attributes are derived state, not durable
--     observations like Records — they don't survive deletion of their Subject.
--
-- Indexes:
--   - subject_id: primary access pattern (BFF reads all attributes for a subject).
--   - type: secondary access pattern (ops queries / cross-subject filtering).

CREATE TABLE attributes (
    id          TEXT PRIMARY KEY,
    subject_id  TEXT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    type        TEXT NOT NULL,
    data        TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (subject_id, type)
);

CREATE INDEX ix_attributes_subject_id ON attributes (subject_id);
CREATE INDEX ix_attributes_type       ON attributes (type);
