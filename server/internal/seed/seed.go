package seed

import (
	"context"
	"time"

	"github.com/open-pm/open-pm/server/internal/api"
	"github.com/open-pm/open-pm/server/internal/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// Run creates the default admin account and workspace if they don't exist.
// This is idempotent — safe to call on every startup.
func Run(ctx context.Context, cfg *config.Config, queries api.Queries) error {
	seed := cfg.Seed

	// Check if admin user already exists
	_, err := queries.GetUserByEmail(ctx, &seed.AdminEmail)
	if err == nil {
		log.Info().Str("email", seed.AdminEmail).Msg("default admin account already exists, skipping seed")
		return nil
	}

	log.Info().Str("email", seed.AdminEmail).Msg("creating default admin account")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seed.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedStr := string(hashedPassword)

	displayName := seed.AdminFirstName
	if seed.AdminLastName != "" {
		displayName += " " + seed.AdminLastName
	}

	now := time.Now()
	user, err := queries.CreateUser(ctx, &seed.AdminEmail, &hashedStr, seed.AdminFirstName, seed.AdminLastName, displayName, &now)
	if err != nil {
		return err
	}

	log.Info().Str("slug", seed.WorkspaceSlug).Msg("creating default workspace")

	workspace, err := queries.CreateWorkspace(ctx, seed.WorkspaceName, seed.WorkspaceSlug, user.ID, nil)
	if err != nil {
		return err
	}

	_, err = queries.CreateWorkspaceMember(ctx, workspace.ID, user.ID, api.RoleOwner)
	if err != nil {
		return err
	}

	log.Info().
		Str("email", seed.AdminEmail).
		Str("password", seed.AdminPassword).
		Str("workspace", seed.WorkspaceSlug).
		Msg("seed complete — default admin account created")

	return nil
}
