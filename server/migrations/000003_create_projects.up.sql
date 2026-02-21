CREATE TABLE projects (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id        UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name                VARCHAR(255) NOT NULL,
    description         TEXT DEFAULT '',
    identifier          VARCHAR(12) NOT NULL,
    network             SMALLINT DEFAULT 2,
    emoji               VARCHAR(255),
    icon_prop           JSONB,
    cover_image_url     TEXT,
    default_assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    project_lead_id     UUID REFERENCES users(id) ON DELETE SET NULL,
    cycle_view          BOOLEAN DEFAULT TRUE,
    module_view         BOOLEAN DEFAULT TRUE,
    page_view           BOOLEAN DEFAULT TRUE,
    inbox_view          BOOLEAN DEFAULT FALSE,
    sort_order          FLOAT DEFAULT 65535,
    archived_at         TIMESTAMPTZ,
    created_by          UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_projects_identifier_workspace ON projects(identifier, workspace_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_projects_name_workspace ON projects(name, workspace_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_workspace_id ON projects(workspace_id);

CREATE TABLE project_members (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role       SMALLINT NOT NULL DEFAULT 15,
    is_active  BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_project_members_unique ON project_members(project_id, user_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_timestamp_projects BEFORE UPDATE ON projects FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_project_members BEFORE UPDATE ON project_members FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
