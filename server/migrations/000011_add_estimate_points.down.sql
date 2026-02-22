DROP TRIGGER IF EXISTS set_timestamp_estimate_systems ON estimate_systems;
DROP INDEX IF EXISTS idx_issues_estimate_point;
DROP INDEX IF EXISTS idx_estimate_systems_project_id;
DROP TABLE IF EXISTS estimate_systems;
ALTER TABLE issues DROP COLUMN IF EXISTS estimate_point;
