package api

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

// Queries defines the database operations used by the API layer.
// This interface will be implemented by the sqlc-generated code wrapped in a repository.
type Queries interface {
	// Users
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email *string) (*User, error)
	CreateUser(ctx context.Context, email *string, encryptedPassword *string, firstName, lastName, displayName string, emailConfirmedAt *time.Time) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, firstName, lastName, displayName, avatarURL *string) (*User, error)
	SetUserLastSignIn(ctx context.Context, id uuid.UUID) error
	SetUserEmailConfirmed(ctx context.Context, id uuid.UUID) error
	SetUserConfirmationToken(ctx context.Context, id uuid.UUID, token string) error
	GetUserByConfirmationToken(ctx context.Context, token string) (*User, error)
	SetUserRecoveryToken(ctx context.Context, id uuid.UUID, token string) error
	GetUserByRecoveryToken(ctx context.Context, token string) (*User, error)

	// Sessions
	CreateSession(ctx context.Context, userID uuid.UUID, ip *string, userAgent *string) (*Session, error)
	DeleteUserSessions(ctx context.Context, userID uuid.UUID) error

	// Refresh Tokens
	CreateRefreshToken(ctx context.Context, sessionID *uuid.UUID, userID uuid.UUID, token string, parent *string) (*RefreshToken, error)
	GetRefreshTokenByToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeRefreshTokenFamily(ctx context.Context, sessionID *uuid.UUID) error

	// Workspaces
	CreateWorkspace(ctx context.Context, name, slug string, ownerID uuid.UUID, organizationSize *string) (*Workspace, error)
	GetWorkspaceByID(ctx context.Context, id uuid.UUID) (*Workspace, error)
	GetWorkspaceBySlug(ctx context.Context, slug string) (*Workspace, error)
	ListWorkspacesByUser(ctx context.Context, userID uuid.UUID) ([]*Workspace, error)
	UpdateWorkspace(ctx context.Context, id uuid.UUID, name, slug, logoURL, organizationSize *string) (*Workspace, error)
	DeleteWorkspace(ctx context.Context, id uuid.UUID) error

	// Workspace Members
	CreateWorkspaceMember(ctx context.Context, workspaceID, userID uuid.UUID, role int16) (*WorkspaceMember, error)
	GetWorkspaceMember(ctx context.Context, workspaceID, userID uuid.UUID) (*WorkspaceMember, error)
	ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID) ([]*WorkspaceMemberWithUser, error)
	UpdateWorkspaceMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role int16) (*WorkspaceMember, error)
	DeleteWorkspaceMember(ctx context.Context, workspaceID, userID uuid.UUID) error

	// Workspace Invites
	CreateWorkspaceInvite(ctx context.Context, workspaceID uuid.UUID, email string, role int16, token string, message *string) (*WorkspaceInvite, error)
	GetWorkspaceInviteByToken(ctx context.Context, token string) (*WorkspaceInvite, error)
	ListWorkspaceInvites(ctx context.Context, workspaceID uuid.UUID) ([]*WorkspaceInvite, error)
	AcceptWorkspaceInvite(ctx context.Context, id uuid.UUID) error

	// Projects
	CreateProject(ctx context.Context, workspaceID uuid.UUID, name, description, identifier string, createdBy, projectLeadID *uuid.UUID) (*Project, error)
	GetProjectByID(ctx context.Context, id uuid.UUID) (*Project, error)
	ListProjectsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*Project, error)
	UpdateProject(ctx context.Context, id uuid.UUID, name, description, identifier, emoji, coverImageURL *string, defaultAssigneeID, projectLeadID *uuid.UUID, network *int16, cycleView, moduleView, pageView, inboxView *bool) (*Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error

	// Project Members
	CreateProjectMember(ctx context.Context, projectID, userID uuid.UUID, role int16) (*ProjectMember, error)
	GetProjectMember(ctx context.Context, projectID, userID uuid.UUID) (*ProjectMember, error)
	ListProjectMembers(ctx context.Context, projectID uuid.UUID) ([]*ProjectMemberWithUser, error)
	UpdateProjectMemberRole(ctx context.Context, projectID, userID uuid.UUID, role int16) (*ProjectMember, error)
	DeleteProjectMember(ctx context.Context, projectID, userID uuid.UUID) error

	// States
	CreateState(ctx context.Context, projectID, workspaceID uuid.UUID, name, description, color, group string, sequence float64, isDefault bool) (*State, error)
	GetStateByID(ctx context.Context, id uuid.UUID) (*State, error)
	ListStatesByProject(ctx context.Context, projectID uuid.UUID) ([]*State, error)
	UpdateState(ctx context.Context, id uuid.UUID, name, description, color *string, group *string, sequence *float64, isDefault *bool) (*State, error)
	DeleteState(ctx context.Context, id uuid.UUID) error

	// Labels
	CreateLabel(ctx context.Context, projectID *uuid.UUID, workspaceID uuid.UUID, parentID *uuid.UUID, name, description, color string, sortOrder float64) (*Label, error)
	GetLabelByID(ctx context.Context, id uuid.UUID) (*Label, error)
	ListLabelsByProject(ctx context.Context, projectID uuid.UUID) ([]*Label, error)
	UpdateLabel(ctx context.Context, id uuid.UUID, name, description, color *string, parentID *uuid.UUID, sortOrder *float64) (*Label, error)
	DeleteLabel(ctx context.Context, id uuid.UUID) error

	// Issues
	CreateIssue(ctx context.Context, params CreateIssueParams) (*Issue, error)
	GetIssueByID(ctx context.Context, id uuid.UUID) (*Issue, error)
	ListIssuesByProject(ctx context.Context, projectID uuid.UUID, filters IssueFilters, limit, offset int) ([]*Issue, error)
	ListSubIssues(ctx context.Context, parentID uuid.UUID) ([]*Issue, error)
	CountIssuesByProject(ctx context.Context, projectID uuid.UUID, filters IssueFilters) (int64, error)
	UpdateIssue(ctx context.Context, id uuid.UUID, params UpdateIssueParams) (*Issue, error)
	DeleteIssue(ctx context.Context, id uuid.UUID) error
	AddIssueAssignee(ctx context.Context, issueID, assigneeID uuid.UUID) error
	RemoveIssueAssignee(ctx context.Context, issueID, assigneeID uuid.UUID) error
	ListIssueAssignees(ctx context.Context, issueID uuid.UUID) ([]*UserSummary, error)
	AddIssueLabel(ctx context.Context, issueID, labelID uuid.UUID) error
	RemoveIssueLabel(ctx context.Context, issueID, labelID uuid.UUID) error
	ListIssueLabels(ctx context.Context, issueID uuid.UUID) ([]*Label, error)

	// Issue Comments
	CreateIssueComment(ctx context.Context, issueID, workspaceID uuid.UUID, commentHTML string, commentJSON []byte, commentStripped string, actorID *uuid.UUID) (*IssueComment, error)
	GetIssueCommentByID(ctx context.Context, id uuid.UUID) (*IssueComment, error)
	ListIssueComments(ctx context.Context, issueID uuid.UUID) ([]*IssueCommentWithUser, error)
	UpdateIssueComment(ctx context.Context, id uuid.UUID, commentHTML string, commentJSON []byte, commentStripped string) (*IssueComment, error)
	DeleteIssueComment(ctx context.Context, id uuid.UUID) error

	// Issue Activities
	CreateIssueActivity(ctx context.Context, params CreateActivityParams) (*IssueActivity, error)
	ListIssueActivities(ctx context.Context, issueID uuid.UUID) ([]*IssueActivityWithUser, error)

	// Work Logs
	CreateWorkLog(ctx context.Context, issueID, projectID, workspaceID uuid.UUID, durationMinutes int, description string, startDate, endDate time.Time, loggedBy uuid.UUID) (*WorkLog, error)
	GetWorkLogByID(ctx context.Context, id uuid.UUID) (*WorkLog, error)
	ListWorkLogsByIssue(ctx context.Context, issueID uuid.UUID) ([]*WorkLogWithUser, error)
	UpdateWorkLog(ctx context.Context, id uuid.UUID, durationMinutes *int, description *string, startDate, endDate *time.Time) (*WorkLog, error)
	DeleteWorkLog(ctx context.Context, id uuid.UUID) error
	SumWorkLogMinutesByIssue(ctx context.Context, issueID uuid.UUID) (int64, error)

	// Cycles
	CreateCycle(ctx context.Context, projectID, workspaceID uuid.UUID, name, description string, startDate, endDate *time.Time, ownedBy uuid.UUID) (*Cycle, error)
	GetCycleByID(ctx context.Context, id uuid.UUID) (*Cycle, error)
	ListCyclesByProject(ctx context.Context, projectID uuid.UUID) ([]*Cycle, error)
	UpdateCycle(ctx context.Context, id uuid.UUID, name, description *string, startDate, endDate *time.Time, sortOrder *float64) (*Cycle, error)
	DeleteCycle(ctx context.Context, id uuid.UUID) error
	AddIssueToCycle(ctx context.Context, cycleID, issueID uuid.UUID) error
	RemoveIssueFromCycle(ctx context.Context, cycleID, issueID uuid.UUID) error
	ListIssuesByCycle(ctx context.Context, cycleID uuid.UUID) ([]*Issue, error)
	CountIssuesByCycle(ctx context.Context, cycleID uuid.UUID) (int64, error)

	// Modules
	CreateModule(ctx context.Context, projectID, workspaceID uuid.UUID, name, description string, startDate, targetDate *time.Time, status string, leadID, createdBy *uuid.UUID) (*Module, error)
	GetModuleByID(ctx context.Context, id uuid.UUID) (*Module, error)
	ListModulesByProject(ctx context.Context, projectID uuid.UUID) ([]*Module, error)
	UpdateModule(ctx context.Context, id uuid.UUID, name, description *string, startDate, targetDate *time.Time, status *string, leadID *uuid.UUID, sortOrder *float64) (*Module, error)
	DeleteModule(ctx context.Context, id uuid.UUID) error
	AddIssueToModule(ctx context.Context, moduleID, issueID uuid.UUID) error
	RemoveIssueFromModule(ctx context.Context, moduleID, issueID uuid.UUID) error
	ListIssuesByModule(ctx context.Context, moduleID uuid.UUID) ([]*Issue, error)
	CountIssuesByModule(ctx context.Context, moduleID uuid.UUID) (int64, error)

	// Pages
	CreatePage(ctx context.Context, projectID *uuid.UUID, workspaceID uuid.UUID, name, descHTML string, descJSON []byte, descStripped string, ownedBy uuid.UUID, parentID *uuid.UUID) (*Page, error)
	GetPageByID(ctx context.Context, id uuid.UUID) (*Page, error)
	ListPagesByProject(ctx context.Context, projectID uuid.UUID) ([]*Page, error)
	UpdatePage(ctx context.Context, id uuid.UUID, name, descHTML *string, descJSON []byte, descStripped *string, color *string, isLocked *bool, parentID *uuid.UUID) (*Page, error)
	DeletePage(ctx context.Context, id uuid.UUID) error

	// Notifications
	CreateNotification(ctx context.Context, workspaceID uuid.UUID, projectID *uuid.UUID, title, message string, data []byte, entityType *string, entityID *uuid.UUID, senderID *uuid.UUID, receiverID uuid.UUID) (*Notification, error)
	ListNotificationsByUser(ctx context.Context, userID, workspaceID uuid.UUID, limit, offset int) ([]*NotificationWithSender, error)
	CountUnreadNotifications(ctx context.Context, userID, workspaceID uuid.UUID) (int64, error)
	MarkNotificationRead(ctx context.Context, id, receiverID uuid.UUID) error
	MarkAllNotificationsRead(ctx context.Context, receiverID, workspaceID uuid.UUID) error

	// Issue Subscribers
	ListIssueSubscribers(ctx context.Context, issueID uuid.UUID) ([]uuid.UUID, error)
	AddIssueSubscriber(ctx context.Context, issueID, subscriberID uuid.UUID) error

	// Estimate Systems
	CreateEstimateSystem(ctx context.Context, projectID, workspaceID uuid.UUID, systemType string, estimates []byte) (*EstimateSystem, error)
	GetEstimateSystemByProject(ctx context.Context, projectID uuid.UUID) (*EstimateSystem, error)
	UpdateEstimateSystem(ctx context.Context, id uuid.UUID, systemType *string, estimates []byte) (*EstimateSystem, error)
	DeleteEstimateSystem(ctx context.Context, id uuid.UUID) error

	// Issue Relations
	AddIssueRelation(ctx context.Context, issueID, relatedIssueID uuid.UUID, relationType string, createdBy *uuid.UUID) (*IssueRelation, error)
	GetIssueRelationByID(ctx context.Context, id uuid.UUID) (*IssueRelation, error)
	RemoveIssueRelation(ctx context.Context, id uuid.UUID) error
	RemoveInverseIssueRelation(ctx context.Context, issueID, relatedIssueID uuid.UUID, relationType string) error
	ListIssueRelations(ctx context.Context, issueID uuid.UUID) ([]*IssueRelationWithIssue, error)

	// Issue Links
	CreateIssueLink(ctx context.Context, issueID uuid.UUID, title, url string, createdBy *uuid.UUID) (*IssueLink, error)
	ListIssueLinksByIssue(ctx context.Context, issueID uuid.UUID) ([]*IssueLink, error)
	UpdateIssueLink(ctx context.Context, id uuid.UUID, title, url *string) (*IssueLink, error)
	DeleteIssueLink(ctx context.Context, id uuid.UUID) error

	// Reports - Issue Analytics
	CountIssuesByState(ctx context.Context, projectID uuid.UUID) ([]*IssuesByStateReport, error)
	CountIssuesByPriority(ctx context.Context, projectID uuid.UUID) ([]*IssuesByPriorityReport, error)
	CountIssuesByType(ctx context.Context, projectID uuid.UUID) ([]*IssuesByTypeReport, error)
	CountIssuesByAssignee(ctx context.Context, projectID uuid.UUID) ([]*IssuesByAssigneeReport, error)

	// Reports - Work Logs
	SumWorkLogsByMember(ctx context.Context, projectID uuid.UUID, startDate, endDate time.Time) ([]*WorkLogSummaryByMember, error)
	SumWorkLogsByIssue(ctx context.Context, projectID uuid.UUID, startDate, endDate time.Time) ([]*WorkLogSummaryByIssue, error)
	SumWorkLogsByDate(ctx context.Context, projectID uuid.UUID, startDate, endDate time.Time) ([]*WorkLogSummaryByDate, error)

	// Reports - Project Health
	ListOverdueIssues(ctx context.Context, projectID uuid.UUID) ([]*Issue, error)
	ListBlockedIssueIDs(ctx context.Context, projectID uuid.UUID) ([]uuid.UUID, error)
	ListUnestimatedIssues(ctx context.Context, projectID uuid.UUID) ([]*Issue, error)

	// Reports - Cycles
	ListCycleIssueIDs(ctx context.Context, cycleID uuid.UUID) ([]uuid.UUID, error)
	CountCycleIssuesByState(ctx context.Context, cycleID uuid.UUID) ([]*IssuesByStateReport, error)
	CountCompletedCycleIssues(ctx context.Context, cycleID uuid.UUID) (int64, error)
	CountTotalCycleIssues(ctx context.Context, cycleID uuid.UUID) (int64, error)

	// Chat Messages
	CreateChatMessage(ctx context.Context, projectID, workspaceID, userID uuid.UUID, role, content string) (*ChatMessage, error)
	ListChatMessages(ctx context.Context, projectID, userID uuid.UUID, limit int) ([]*ChatMessage, error)
	DeleteChatHistory(ctx context.Context, projectID, userID uuid.UUID) error
	// Search
	SearchIssues(ctx context.Context, workspaceID uuid.UUID, query string, limit int) ([]*Issue, error)
	SearchPages(ctx context.Context, workspaceID uuid.UUID, query string, limit int) ([]*Page, error)

	// File Assets
	CreateFileAsset(ctx context.Context, workspaceID uuid.UUID, entityType string, entityID uuid.UUID, fileName string, fileSize int64, contentType, storageKey string, uploadedBy uuid.UUID) (*FileAsset, error)
	GetFileAssetByID(ctx context.Context, id uuid.UUID) (*FileAsset, error)
	ListFileAssetsByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*FileAsset, error)
	DeleteFileAsset(ctx context.Context, id uuid.UUID) error
}
