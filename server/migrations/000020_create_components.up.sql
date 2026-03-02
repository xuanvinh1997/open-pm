-- Components: project sub-sections with optional leads and default-assignee routing
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'component_default_assignee') THEN
        CREATE TYPE component_default_assignee AS ENUM (
            'project_default', 'component_lead', 'unassigned'
        );
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS components (
    id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id            UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id          UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name                  VARCHAR(255) NOT NULL,
    description           TEXT DEFAULT '',
    lead_id               UUID REFERENCES users(id) ON DELETE SET NULL,
    default_assignee_type component_default_assignee DEFAULT 'project_default',
    sort_order            FLOAT DEFAULT 65535,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ,
    UNIQUE(name, project_id)
);

CREATE INDEX IF NOT EXISTS idx_components_project_id ON components(project_id) WHERE deleted_at IS NULL;

-- Junction table: which components an issue belongs to
CREATE TABLE IF NOT EXISTS issue_components (
    issue_id     UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES components(id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (issue_id, component_id)
);

CREATE INDEX IF NOT EXISTS idx_issue_components_issue_id ON issue_components(issue_id);
CREATE INDEX IF NOT EXISTS idx_issue_components_component_id ON issue_components(component_id);

CREATE TRIGGER set_timestamp_components BEFORE UPDATE ON components FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
