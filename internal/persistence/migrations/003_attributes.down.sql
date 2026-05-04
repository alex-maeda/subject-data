-- Migration 003 down: drop the attributes table.

DROP INDEX IF EXISTS ix_attributes_type;
DROP TABLE IF EXISTS attributes;
