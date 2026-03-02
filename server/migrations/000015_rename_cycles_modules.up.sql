-- Rename cycle_issues column before renaming tables
ALTER TABLE cycle_issues RENAME COLUMN cycle_id TO sprint_id;
ALTER TABLE module_issues RENAME COLUMN module_id TO epic_id;
ALTER TABLE module_members RENAME COLUMN module_id TO epic_id;

-- Drop triggers (must drop before renaming tables they reference)
DROP TRIGGER IF EXISTS set_timestamp_cycle_issues ON cycle_issues;
DROP TRIGGER IF EXISTS set_timestamp_cycles ON cycles;
DROP TRIGGER IF EXISTS set_timestamp_module_issues ON module_issues;
DROP TRIGGER IF EXISTS set_timestamp_modules ON modules;

-- Rename tables
ALTER TABLE cycles RENAME TO sprints;
ALTER TABLE cycle_issues RENAME TO sprint_issues;
ALTER TABLE modules RENAME TO epics;
ALTER TABLE module_issues RENAME TO epic_issues;
ALTER TABLE module_members RENAME TO epic_members;

-- Rename indexes
ALTER INDEX idx_cycles_project_id RENAME TO idx_sprints_project_id;
ALTER INDEX idx_cycle_issues_unique RENAME TO idx_sprint_issues_unique;
ALTER INDEX idx_modules_project_id RENAME TO idx_epics_project_id;
ALTER INDEX idx_module_issues_unique RENAME TO idx_epic_issues_unique;
ALTER INDEX idx_module_members_unique RENAME TO idx_epic_members_unique;

-- Recreate triggers with new names on renamed tables
CREATE TRIGGER set_timestamp_sprints BEFORE UPDATE ON sprints FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_sprint_issues BEFORE UPDATE ON sprint_issues FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_epics BEFORE UPDATE ON epics FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_epic_issues BEFORE UPDATE ON epic_issues FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();

-- Rename project columns
ALTER TABLE projects RENAME COLUMN cycle_view TO sprint_view;
ALTER TABLE projects RENAME COLUMN module_view TO epic_view;
