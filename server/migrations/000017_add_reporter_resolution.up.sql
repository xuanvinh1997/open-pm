-- Add reporter_id to issues (separate from created_by, independently updatable)
ALTER TABLE issues ADD COLUMN IF NOT EXISTS reporter_id UUID REFERENCES users(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_issues_reporter_id ON issues(reporter_id) WHERE deleted_at IS NULL;

-- Set reporter_id = created_by for all existing issues
UPDATE issues SET reporter_id = created_by WHERE reporter_id IS NULL AND created_by IS NOT NULL;

-- Issue resolutions table (Fixed, Won't Fix, Duplicate, Cannot Reproduce, etc.)
CREATE TABLE IF NOT EXISTS issue_resolutions (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    description  TEXT DEFAULT '',
    is_default   BOOLEAN DEFAULT FALSE,
    sort_order   FLOAT DEFAULT 65535,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(name, project_id)
);

CREATE INDEX IF NOT EXISTS idx_issue_resolutions_project_id ON issue_resolutions(project_id);

-- Add resolution tracking to issues
ALTER TABLE issues ADD COLUMN IF NOT EXISTS resolution_id UUID REFERENCES issue_resolutions(id) ON DELETE SET NULL;
ALTER TABLE issues ADD COLUMN IF NOT EXISTS resolved_at TIMESTAMPTZ;
