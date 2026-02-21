DROP TRIGGER IF EXISTS set_timestamp_workspace_invites ON workspace_member_invites;
DROP TRIGGER IF EXISTS set_timestamp_workspace_members ON workspace_members;
DROP TRIGGER IF EXISTS set_timestamp_workspaces ON workspaces;
DROP TABLE IF EXISTS workspace_member_invites;
DROP TABLE IF EXISTS workspace_members;
DROP TABLE IF EXISTS workspaces;
