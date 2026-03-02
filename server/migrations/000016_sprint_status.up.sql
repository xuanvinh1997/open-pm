ALTER TABLE sprints ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'planned'
  CHECK (status IN ('planned', 'active', 'completed'));
