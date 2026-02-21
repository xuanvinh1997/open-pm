CREATE TYPE issue_type AS ENUM ('story', 'bug', 'task', 'epic');
ALTER TABLE issues ADD COLUMN issue_type issue_type DEFAULT 'task';
CREATE INDEX idx_issues_issue_type ON issues(issue_type);
