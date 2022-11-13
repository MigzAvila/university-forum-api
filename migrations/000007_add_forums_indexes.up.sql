-- Filename: migrations/000005_add_forums_indexes.up.sql
CREATE INDEX IF NOT EXISTS forums_title_idx ON forums USING GIN(to_tsvector('simple', title));