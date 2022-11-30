-- Create schema
CREATE SCHEMA taxi;
SET search_path TO taxi;

-- Set wal level
ALTER SYSTEM SET wal_level = logical;