-- Add sprint goal field for Scrum ceremony alignment
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS goal TEXT DEFAULT '';
