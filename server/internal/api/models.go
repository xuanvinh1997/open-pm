package api

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

// Domain models used by the API layer.
// These mirror the database tables and will be populated by the repository layer.

type User struct {
	ID                 uuid.UUID        `json:"id"`
	Email              *string          `json:"email"`
	EncryptedPassword  *string          `json:"-"`
	EmailConfirmedAt   *time.Time       `json:"email_confirmed_at,omitempty"`
	FirstName          string           `json:"first_name"`
	LastName           string           `json:"last_name"`
	AvatarURL          string           `json:"avatar_url"`
	DisplayName        string           `json:"display_name"`
	IsActive           bool             `json:"is_active"`
	UserMetadata       json.RawMessage  `json:"user_metadata,omitempty"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}

type UserSummary struct {
	ID          uuid.UUID `json:"id"`
	Email       *string   `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
}

type Session struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	IP        *string    `json:"ip,omitempty"`
	UserAgent *string    `json:"user_agent,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	NotAfter  *time.Time `json:"not_after,omitempty"`
}

type RefreshToken struct {
	ID        int64      `json:"id"`
	SessionID *uuid.UUID `json:"session_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Token     string     `json:"-"`
	Parent    *string    `json:"-"`
	Revoked   bool       `json:"revoked"`
	CreatedAt time.Time  `json:"created_at"`
}

type Workspace struct {
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	LogoURL          string     `json:"logo_url"`
	OwnerID          uuid.UUID  `json:"owner_id"`
	OrganizationSize *string    `json:"organization_size,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type WorkspaceMember struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	UserID      uuid.UUID `json:"user_id"`
	Role        int16     `json:"role"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type WorkspaceMemberWithUser struct {
	WorkspaceMember
	Email       *string `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DisplayName string  `json:"display_name"`
	AvatarURL   string  `json:"avatar_url"`
}

type WorkspaceInvite struct {
	ID          uuid.UUID  `json:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	Email       string     `json:"email"`
	Role        int16      `json:"role"`
	Token       string     `json:"-"`
	Message     *string    `json:"message,omitempty"`
	Accepted    bool       `json:"accepted"`
	RespondedAt *time.Time `json:"responded_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Project struct {
	ID                uuid.UUID  `json:"id"`
	WorkspaceID       uuid.UUID  `json:"workspace_id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	Identifier        string     `json:"identifier"`
	Network           int16      `json:"network"`
	Emoji             *string    `json:"emoji,omitempty"`
	CoverImageURL     *string    `json:"cover_image_url,omitempty"`
	DefaultAssigneeID *uuid.UUID `json:"default_assignee_id,omitempty"`
	ProjectLeadID     *uuid.UUID `json:"project_lead_id,omitempty"`
	CycleView         bool       `json:"cycle_view"`
	ModuleView        bool       `json:"module_view"`
	PageView          bool       `json:"page_view"`
	InboxView         bool       `json:"inbox_view"`
	SortOrder         float64    `json:"sort_order"`
	ArchivedAt        *time.Time `json:"archived_at,omitempty"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type ProjectMember struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      int16     `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ProjectMemberWithUser struct {
	ProjectMember
	Email       *string `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DisplayName string  `json:"display_name"`
	AvatarURL   string  `json:"avatar_url"`
}

type State struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"project_id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Group       string    `json:"group"`
	Sequence    float64   `json:"sequence"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Label struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   *uuid.UUID `json:"project_id,omitempty"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Color       string     `json:"color"`
	SortOrder   float64    `json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Issue struct {
	ID                 uuid.UUID       `json:"id"`
	ProjectID          uuid.UUID       `json:"project_id"`
	WorkspaceID        uuid.UUID       `json:"workspace_id"`
	ParentID           *uuid.UUID      `json:"parent_id,omitempty"`
	StateID            *uuid.UUID      `json:"state_id,omitempty"`
	Name               string          `json:"name"`
	DescriptionHTML    string          `json:"description_html"`
	DescriptionJSON    json.RawMessage `json:"description_json,omitempty"`
	DescriptionStripped string         `json:"description_stripped"`
	Priority           string          `json:"priority"`
	IssueType          string          `json:"issue_type"`
	StartDate          *time.Time      `json:"start_date,omitempty"`
	TargetDate         *time.Time      `json:"target_date,omitempty"`
	SequenceID         int             `json:"sequence_id"`
	SortOrder          float64         `json:"sort_order"`
	CompletedAt        *time.Time      `json:"completed_at,omitempty"`
	ArchivedAt         *time.Time      `json:"archived_at,omitempty"`
	IsDraft            bool            `json:"is_draft"`
	CreatedBy          *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID      `json:"updated_by,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type CreateIssueParams struct {
	ProjectID          uuid.UUID
	WorkspaceID        uuid.UUID
	ParentID           *uuid.UUID
	StateID            *uuid.UUID
	Name               string
	DescriptionHTML    string
	DescriptionJSON    json.RawMessage
	DescriptionStripped string
	Priority           string
	IssueType          string
	StartDate          *time.Time
	TargetDate         *time.Time
	IsDraft            bool
	CreatedBy          *uuid.UUID
}

type UpdateIssueParams struct {
	StateID            *uuid.UUID
	Name               *string
	DescriptionHTML    *string
	DescriptionJSON    json.RawMessage
	DescriptionStripped *string
	Priority           *string
	IssueType          *string
	StartDate          *time.Time
	TargetDate         *time.Time
	ParentID           *uuid.UUID
	SortOrder          *float64
	IsDraft            *bool
	CompletedAt        *time.Time
	ArchivedAt         *time.Time
	UpdatedBy          *uuid.UUID
}

type IssueComment struct {
	ID              uuid.UUID       `json:"id"`
	IssueID         uuid.UUID       `json:"issue_id"`
	WorkspaceID     uuid.UUID       `json:"workspace_id"`
	CommentHTML     string          `json:"comment_html"`
	CommentJSON     json.RawMessage `json:"comment_json,omitempty"`
	CommentStripped string          `json:"comment_stripped"`
	ActorID         *uuid.UUID      `json:"actor_id,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type IssueCommentWithUser struct {
	IssueComment
	Email       *string `json:"email,omitempty"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DisplayName string  `json:"display_name"`
	AvatarURL   string  `json:"avatar_url"`
}

type IssueActivity struct {
	ID            uuid.UUID  `json:"id"`
	IssueID       *uuid.UUID `json:"issue_id,omitempty"`
	ProjectID     uuid.UUID  `json:"project_id"`
	WorkspaceID   uuid.UUID  `json:"workspace_id"`
	Verb          string     `json:"verb"`
	Field         *string    `json:"field,omitempty"`
	OldValue      *string    `json:"old_value,omitempty"`
	NewValue      *string    `json:"new_value,omitempty"`
	OldIdentifier *uuid.UUID `json:"old_identifier,omitempty"`
	NewIdentifier *uuid.UUID `json:"new_identifier,omitempty"`
	Comment       string     `json:"comment"`
	ActorID       *uuid.UUID `json:"actor_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type IssueActivityWithUser struct {
	IssueActivity
	Email       *string `json:"email,omitempty"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DisplayName string  `json:"display_name"`
	AvatarURL   string  `json:"avatar_url"`
}

type CreateActivityParams struct {
	IssueID       *uuid.UUID
	ProjectID     uuid.UUID
	WorkspaceID   uuid.UUID
	Verb          string
	Field         *string
	OldValue      *string
	NewValue      *string
	OldIdentifier *uuid.UUID
	NewIdentifier *uuid.UUID
	Comment       string
	ActorID       *uuid.UUID
}

type Cycle struct {
	ID               uuid.UUID       `json:"id"`
	ProjectID        uuid.UUID       `json:"project_id"`
	WorkspaceID      uuid.UUID       `json:"workspace_id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	StartDate        *time.Time      `json:"start_date,omitempty"`
	EndDate          *time.Time      `json:"end_date,omitempty"`
	OwnedBy          uuid.UUID       `json:"owned_by"`
	SortOrder        float64         `json:"sort_order"`
	ProgressSnapshot json.RawMessage `json:"progress_snapshot,omitempty"`
	ArchivedAt       *time.Time      `json:"archived_at,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type Module struct {
	ID              uuid.UUID  `json:"id"`
	ProjectID       uuid.UUID  `json:"project_id"`
	WorkspaceID     uuid.UUID  `json:"workspace_id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	DescriptionHTML string     `json:"description_html"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	TargetDate      *time.Time `json:"target_date,omitempty"`
	Status          string     `json:"status"`
	LeadID          *uuid.UUID `json:"lead_id,omitempty"`
	SortOrder       float64    `json:"sort_order"`
	ArchivedAt      *time.Time `json:"archived_at,omitempty"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Page struct {
	ID                 uuid.UUID       `json:"id"`
	ProjectID          *uuid.UUID      `json:"project_id,omitempty"`
	WorkspaceID        uuid.UUID       `json:"workspace_id"`
	Name               string          `json:"name"`
	DescriptionHTML    string          `json:"description_html"`
	DescriptionJSON    json.RawMessage `json:"description_json,omitempty"`
	DescriptionStripped string         `json:"description_stripped"`
	Color              string          `json:"color"`
	IsLocked           bool            `json:"is_locked"`
	ArchivedAt         *time.Time      `json:"archived_at,omitempty"`
	OwnedBy            uuid.UUID       `json:"owned_by"`
	ParentID           *uuid.UUID      `json:"parent_id,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type WorkLog struct {
	ID              uuid.UUID `json:"id"`
	IssueID         uuid.UUID `json:"issue_id"`
	ProjectID       uuid.UUID `json:"project_id"`
	WorkspaceID     uuid.UUID `json:"workspace_id"`
	DurationMinutes int       `json:"duration_minutes"`
	Description     string    `json:"description"`
	LoggedAt        time.Time `json:"logged_at"`
	LoggedBy        uuid.UUID `json:"logged_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type WorkLogWithUser struct {
	WorkLog
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type Notification struct {
	ID          uuid.UUID       `json:"id"`
	WorkspaceID uuid.UUID       `json:"workspace_id"`
	ProjectID   *uuid.UUID      `json:"project_id,omitempty"`
	Title       string          `json:"title"`
	Message     string          `json:"message"`
	Data        json.RawMessage `json:"data,omitempty"`
	EntityType  *string         `json:"entity_type,omitempty"`
	EntityID    *uuid.UUID      `json:"entity_id,omitempty"`
	SenderID    *uuid.UUID      `json:"sender_id,omitempty"`
	ReceiverID  uuid.UUID       `json:"receiver_id"`
	ReadAt      *time.Time      `json:"read_at,omitempty"`
	SnoozedTill *time.Time      `json:"snoozed_till,omitempty"`
	ArchivedAt  *time.Time      `json:"archived_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}
