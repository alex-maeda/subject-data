-- Migration 003 down: drop the attributes table.

DROP INDEX IF EXISTS ix_attributes_type;
DROP INDEX IF EXISTS ix_attributes_subject_id;
DROP TABLE IF EXISTS attributes;
