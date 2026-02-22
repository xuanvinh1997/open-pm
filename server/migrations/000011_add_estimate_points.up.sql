-- Add estimate_point column to issues
ALTER TABLE issues ADD COLUMN estimate_point INTEGER;

-- Estimate systems: per-project configuration for estimate scales
CREATE TABLE estimate_systems (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    type         VARCHAR(50) NOT NULL DEFAULT 'points',
    estimates    JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_estimate_systems_project_id ON estimate_systems(project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_issues_estimate_point ON issues(estimate_point) WHERE deleted_at IS NULL AND estimate_point IS NOT NULL;

CREATE TRIGGER set_timestamp_estimate_systems BEFORE UPDATE ON estimate_systems FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
