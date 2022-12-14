 -- filename migrations/000001_add_user_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    create_at timestamp(0) without time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL
)