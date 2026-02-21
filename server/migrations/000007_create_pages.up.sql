CREATE TABLE pages (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id           UUID REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id         UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name                 VARCHAR(255) NOT NULL,
    description_html     TEXT DEFAULT '<p></p>',
    description_json     JSONB DEFAULT '{}'::jsonb,
    description_stripped TEXT DEFAULT '',
    color                VARCHAR(255) DEFAULT '',
    is_locked            BOOLEAN DEFAULT FALSE,
    archived_at          TIMESTAMPTZ,
    owned_by             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id            UUID REFERENCES pages(id) ON DELETE SET NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ
);

CREATE INDEX idx_pages_project_id ON pages(project_id);
CREATE INDEX idx_pages_workspace_id ON pages(workspace_id);
CREATE INDEX idx_pages_owned_by ON pages(owned_by);

CREATE TABLE page_labels (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id    UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    label_id   UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_page_labels_unique ON page_labels(page_id, label_id);

CREATE TRIGGER set_timestamp_pages BEFORE UPDATE ON pages FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
