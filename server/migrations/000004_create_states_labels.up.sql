CREATE TYPE state_group AS ENUM ('backlog', 'unstarted', 'started', 'completed', 'cancelled', 'triage');

CREATE TABLE states (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    description  TEXT DEFAULT '',
    color        VARCHAR(255) NOT NULL DEFAULT '#60646C',
    "group"      state_group NOT NULL DEFAULT 'backlog',
    sequence     FLOAT DEFAULT 65535,
    is_default   BOOLEAN DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_states_name_project ON states(name, project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_states_project_id ON states(project_id);

CREATE TABLE labels (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id   UUID REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    parent_id    UUID REFERENCES labels(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    description  TEXT DEFAULT '',
    color        VARCHAR(255) NOT NULL DEFAULT '#60646C',
    sort_order   FLOAT DEFAULT 65535,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_labels_project_id ON labels(project_id);
CREATE INDEX idx_labels_workspace_id ON labels(workspace_id);

CREATE TRIGGER set_timestamp_states BEFORE UPDATE ON states FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_labels BEFORE UPDATE ON labels FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
