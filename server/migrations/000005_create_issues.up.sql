CREATE TYPE issue_priority AS ENUM ('urgent', 'high', 'medium', 'low', 'none');

CREATE TABLE issues (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id           UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id         UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    parent_id            UUID REFERENCES issues(id) ON DELETE CASCADE,
    state_id             UUID REFERENCES states(id) ON DELETE SET NULL,
    name                 VARCHAR(255) NOT NULL,
    description_html     TEXT DEFAULT '<p></p>',
    description_json     JSONB DEFAULT '{}'::jsonb,
    description_stripped TEXT DEFAULT '',
    priority             issue_priority DEFAULT 'none',
    start_date           DATE,
    target_date          DATE,
    sequence_id          INT NOT NULL DEFAULT 0,
    sort_order           FLOAT DEFAULT 65535,
    completed_at         TIMESTAMPTZ,
    archived_at          TIMESTAMPTZ,
    is_draft             BOOLEAN DEFAULT FALSE,
    created_by           UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by           UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ
);

CREATE INDEX idx_issues_project_id ON issues(project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_issues_workspace_id ON issues(workspace_id);
CREATE INDEX idx_issues_state_id ON issues(state_id);
CREATE INDEX idx_issues_parent_id ON issues(parent_id);
CREATE INDEX idx_issues_priority ON issues(priority);
CREATE INDEX idx_issues_sequence ON issues(project_id, sequence_id);

-- Auto-increment sequence_id per project
CREATE OR REPLACE FUNCTION set_issue_sequence_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.sequence_id := COALESCE(
        (SELECT MAX(sequence_id) + 1 FROM issues WHERE project_id = NEW.project_id),
        1
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_issue_sequence_id
    BEFORE INSERT ON issues FOR EACH ROW
    EXECUTE FUNCTION set_issue_sequence_id();

-- Issue assignees
CREATE TABLE issue_assignees (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    issue_id    UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    assignee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_issue_assignees_unique ON issue_assignees(issue_id, assignee_id);

-- Issue labels
CREATE TABLE issue_labels (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    label_id   UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_issue_labels_unique ON issue_labels(issue_id, label_id);

-- Issue comments
CREATE TABLE issue_comments (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    issue_id         UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    workspace_id     UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    comment_html     TEXT DEFAULT '',
    comment_json     JSONB DEFAULT '{}'::jsonb,
    comment_stripped  TEXT DEFAULT '',
    actor_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_issue_comments_issue_id ON issue_comments(issue_id);

-- Issue activity log
CREATE TABLE issue_activities (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    issue_id       UUID REFERENCES issues(id) ON DELETE CASCADE,
    project_id     UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    workspace_id   UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    verb           VARCHAR(255) NOT NULL DEFAULT 'created',
    field          VARCHAR(255),
    old_value      TEXT,
    new_value      TEXT,
    old_identifier UUID,
    new_identifier UUID,
    comment        TEXT DEFAULT '',
    actor_id       UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_issue_activities_issue_id ON issue_activities(issue_id);

-- Issue subscribers
CREATE TABLE issue_subscribers (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    issue_id      UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    subscriber_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_issue_subscribers_unique ON issue_subscribers(issue_id, subscriber_id);

CREATE TRIGGER set_timestamp_issues BEFORE UPDATE ON issues FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_issue_comments BEFORE UPDATE ON issue_comments FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
