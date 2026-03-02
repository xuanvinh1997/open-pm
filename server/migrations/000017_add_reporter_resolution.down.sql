ALTER TABLE issues DROP COLUMN IF EXISTS resolved_at;
ALTER TABLE issues DROP COLUMN IF EXISTS resolution_id;
DROP TABLE IF EXISTS issue_resolutions;
DROP INDEX IF EXISTS idx_issues_reporter_id;
ALTER TABLE issues DROP COLUMN IF EXISTS reporter_id;
