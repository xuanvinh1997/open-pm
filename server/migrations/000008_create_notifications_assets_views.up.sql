CREATE TABLE notifications (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    project_id   UUID REFERENCES projects(id) ON DELETE CASCADE,
    title        VARCHAR(255) NOT NULL,
    message      TEXT DEFAULT '',
    data         JSONB DEFAULT '{}'::jsonb,
    entity_type  VARCHAR(100),
    entity_id    UUID,
    sender_id    UUID REFERENCES users(id) ON DELETE SET NULL,
    receiver_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    read_at      TIMESTAMPTZ,
    snoozed_till TIMESTAMPTZ,
    archived_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_receiver_id ON notifications(receiver_id);
CREATE INDEX idx_notifications_workspace_id ON notifications(workspace_id);
CREATE INDEX idx_notifications_unread ON notifications(receiver_id) WHERE read_at IS NULL;

CREATE TABLE file_assets (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    asset_name   VARCHAR(255) NOT NULL,
    asset_key    TEXT NOT NULL,
    asset_url    TEXT NOT NULL,
    asset_size   BIGINT DEFAULT 0,
    content_type VARCHAR(255) DEFAULT '',
    entity_type  VARCHAR(100),
    entity_id    UUID,
    uploaded_by  UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_file_assets_workspace_id ON file_assets(workspace_id);
CREATE INDEX idx_file_assets_entity ON file_assets(entity_type, entity_id);

CREATE TABLE views (
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id         UUID REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id       UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name               VARCHAR(255) NOT NULL,
    description        TEXT DEFAULT '',
    filters            JSONB DEFAULT '{}'::jsonb,
    display_filters    JSONB DEFAULT '{}'::jsonb,
    display_properties JSONB DEFAULT '{}'::jsonb,
    access             SMALLINT DEFAULT 1,
    sort_order         FLOAT DEFAULT 65535,
    created_by         UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);

CREATE INDEX idx_views_project_id ON views(project_id);
CREATE INDEX idx_views_workspace_id ON views(workspace_id);

CREATE TRIGGER set_timestamp_notifications BEFORE UPDATE ON notifications FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_file_assets BEFORE UPDATE ON file_assets FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_views BEFORE UPDATE ON views FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
