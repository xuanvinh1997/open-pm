CREATE TABLE workspaces (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name              VARCHAR(80) NOT NULL,
    slug              VARCHAR(48) NOT NULL,
    logo_url          TEXT DEFAULT '',
    owner_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_size VARCHAR(20),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_workspaces_slug ON workspaces(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_workspaces_owner_id ON workspaces(owner_id);

CREATE TABLE workspace_members (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role         SMALLINT NOT NULL DEFAULT 15,
    is_active    BOOLEAN DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_workspace_members_unique ON workspace_members(workspace_id, user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_workspace_members_user_id ON workspace_members(user_id);

CREATE TABLE workspace_member_invites (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    email        VARCHAR(255) NOT NULL,
    role         SMALLINT NOT NULL DEFAULT 15,
    token        VARCHAR(255) NOT NULL,
    message      TEXT,
    accepted     BOOLEAN DEFAULT FALSE,
    responded_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_workspace_invites_unique ON workspace_member_invites(email, workspace_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_timestamp_workspaces BEFORE UPDATE ON workspaces FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_workspace_members BEFORE UPDATE ON workspace_members FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_workspace_invites BEFORE UPDATE ON workspace_member_invites FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
