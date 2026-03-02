-- Versions / Releases: track fix versions and affected versions on issues
CREATE TABLE IF NOT EXISTS versions (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    description  TEXT DEFAULT '',
    start_date   DATE,
    release_date DATE,
    released     BOOLEAN DEFAULT FALSE,
    released_at  TIMESTAMPTZ,
    archived_at  TIMESTAMPTZ,
    sort_order   FLOAT DEFAULT 65535,
    created_by   UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,
    UNIQUE(name, project_id)
);

CREATE INDEX IF NOT EXISTS idx_versions_project_id ON versions(project_id) WHERE deleted_at IS NULL;

-- Fix versions: which version an issue is fixed in
CREATE TABLE IF NOT EXISTS issue_fix_versions (
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    version_id UUID NOT NULL REFERENCES versions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (issue_id, version_id)
);

-- Affects versions: which versions an issue affects
CREATE TABLE IF NOT EXISTS issue_affects_versions (
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    version_id UUID NOT NULL REFERENCES versions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (issue_id, version_id)
);

CREATE TRIGGER set_timestamp_versions BEFORE UPDATE ON versions FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
