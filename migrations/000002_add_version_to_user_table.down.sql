 -- filename migrations/000002_add_version_to_user_table.down.sql
  
ALTER TABLE users 
DROP COLUMN IF EXISTS version;