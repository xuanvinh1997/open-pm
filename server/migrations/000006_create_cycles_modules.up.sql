CREATE TABLE cycles (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id        UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id      UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name              VARCHAR(255) NOT NULL,
    description       TEXT DEFAULT '',
    start_date        TIMESTAMPTZ,
    end_date          TIMESTAMPTZ,
    owned_by          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    view_props        JSONB DEFAULT '{}'::jsonb,
    sort_order        FLOAT DEFAULT 65535,
    progress_snapshot JSONB DEFAULT '{}'::jsonb,
    archived_at       TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_cycles_project_id ON cycles(project_id);

CREATE TABLE cycle_issues (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cycle_id   UUID NOT NULL REFERENCES cycles(id) ON DELETE CASCADE,
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_cycle_issues_unique ON cycle_issues(cycle_id, issue_id);

CREATE TABLE modules (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id       UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id     UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name             VARCHAR(255) NOT NULL,
    description      TEXT DEFAULT '',
    description_html TEXT DEFAULT '',
    description_json JSONB DEFAULT '{}'::jsonb,
    start_date       DATE,
    target_date      DATE,
    status           VARCHAR(50) DEFAULT 'backlog',
    lead_id          UUID REFERENCES users(id) ON DELETE SET NULL,
    sort_order       FLOAT DEFAULT 65535,
    archived_at      TIMESTAMPTZ,
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_modules_project_id ON modules(project_id);

CREATE TABLE module_issues (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    module_id  UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_module_issues_unique ON module_issues(module_id, issue_id);

CREATE TABLE module_members (
    id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role      VARCHAR(50) DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_module_members_unique ON module_members(module_id, member_id);

CREATE TRIGGER set_timestamp_cycles BEFORE UPDATE ON cycles FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_cycle_issues BEFORE UPDATE ON cycle_issues FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_modules BEFORE UPDATE ON modules FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_module_issues BEFORE UPDATE ON module_issues FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
