CREATE TABLE issue_links (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    url        TEXT NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_issue_links_issue_id ON issue_links(issue_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_timestamp_issue_links BEFORE UPDATE ON issue_links FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
