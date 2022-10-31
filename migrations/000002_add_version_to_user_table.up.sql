 -- filename 000001_add_version_to_user_table.up.sql
 
ALTER TABLE users 
ADD COLUMN version integer NOT NULL DEFAULT 1;

