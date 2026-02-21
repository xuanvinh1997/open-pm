DROP TRIGGER IF EXISTS set_timestamp_refresh_tokens ON refresh_tokens;
DROP TRIGGER IF EXISTS set_timestamp_sessions ON sessions;
DROP TRIGGER IF EXISTS set_timestamp_users ON users;
DROP FUNCTION IF EXISTS trigger_set_timestamp();
DROP TABLE IF EXISTS one_time_tokens;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS identities;
DROP TABLE IF EXISTS users;
