CREATE TABLE work_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    issue_id UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    duration_minutes INT NOT NULL CHECK (duration_minutes > 0),
    description TEXT NOT NULL DEFAULT '',
    logged_at DATE NOT NULL DEFAULT CURRENT_DATE,
    logged_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_work_logs_issue_id ON work_logs(issue_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_work_logs_logged_by ON work_logs(logged_by) WHERE deleted_at IS NULL;

CREATE TRIGGER set_timestamp_work_logs BEFORE UPDATE ON work_logs FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
