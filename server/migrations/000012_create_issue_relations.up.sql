CREATE TYPE relation_type AS ENUM ('relates_to', 'blocks', 'blocked_by', 'duplicate_of');

CREATE TABLE issue_relations (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    issue_id         UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    related_issue_id UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    relation_type    relation_type NOT NULL DEFAULT 'relates_to',
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT issue_relations_no_self CHECK (issue_id != related_issue_id)
);

CREATE UNIQUE INDEX idx_issue_relations_unique ON issue_relations(issue_id, related_issue_id, relation_type);
CREATE INDEX idx_issue_relations_issue_id ON issue_relations(issue_id);
CREATE INDEX idx_issue_relations_related_issue_id ON issue_relations(related_issue_id);
