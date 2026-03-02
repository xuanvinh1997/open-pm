-- Reverse project columns
ALTER TABLE projects RENAME COLUMN sprint_view TO cycle_view;
ALTER TABLE projects RENAME COLUMN epic_view TO module_view;

-- Drop triggers
DROP TRIGGER IF EXISTS set_timestamp_sprint_issues ON sprint_issues;
DROP TRIGGER IF EXISTS set_timestamp_sprints ON sprints;
DROP TRIGGER IF EXISTS set_timestamp_epic_issues ON epic_issues;
DROP TRIGGER IF EXISTS set_timestamp_epics ON epics;

-- Rename tables back
ALTER TABLE sprints RENAME TO cycles;
ALTER TABLE sprint_issues RENAME TO cycle_issues;
ALTER TABLE epics RENAME TO modules;
ALTER TABLE epic_issues RENAME TO module_issues;
ALTER TABLE epic_members RENAME TO module_members;

-- Rename columns back
ALTER TABLE cycle_issues RENAME COLUMN sprint_id TO cycle_id;
ALTER TABLE module_issues RENAME COLUMN epic_id TO module_id;
ALTER TABLE module_members RENAME COLUMN epic_id TO module_id;

-- Rename indexes back
ALTER INDEX idx_sprints_project_id RENAME TO idx_cycles_project_id;
ALTER INDEX idx_sprint_issues_unique RENAME TO idx_cycle_issues_unique;
ALTER INDEX idx_epics_project_id RENAME TO idx_modules_project_id;
ALTER INDEX idx_epic_issues_unique RENAME TO idx_module_issues_unique;
ALTER INDEX idx_epic_members_unique RENAME TO idx_module_members_unique;

-- Recreate triggers with original names
CREATE TRIGGER set_timestamp_cycles BEFORE UPDATE ON cycles FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_cycle_issues BEFORE UPDATE ON cycle_issues FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_modules BEFORE UPDATE ON modules FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_module_issues BEFORE UPDATE ON module_issues FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
