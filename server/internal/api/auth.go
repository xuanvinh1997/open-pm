package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int         `json:"expires_in"`
	User         interface{} `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RecoverRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" validate:"required,min=8"`
}

// Signup handles POST /auth/signup
func (a *API) Signup(w http.ResponseWriter, r *http.Request) error {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	// Check if user already exists
	_, err := a.queries.GetUserByEmail(r.Context(), &req.Email)
	if err == nil {
		return conflictError("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return internalServerError("failed to hash password")
	}

	displayName := req.FirstName
	if req.LastName != "" {
		displayName += " " + req.LastName
	}

	hashedStr := string(hashedPassword)

	// In development, auto-confirm email
	var confirmedAt *time.Time
	if a.config.Environment == "development" {
		now := time.Now()
		confirmedAt = &now
	}

	user, err := a.queries.CreateUser(r.Context(), &req.Email, &hashedStr, req.FirstName, req.LastName, displayName, confirmedAt)
	if err != nil {
		return internalServerError("failed to create user")
	}

	// Create session and tokens
	tokenResp, err := a.createSession(r, user.ID)
	if err != nil {
		return err
	}
	tokenResp.User = sanitizeUser(user)

	return sendJSON(w, http.StatusCreated, tokenResp)
}

// Token handles POST /auth/token
func (a *API) Token(w http.ResponseWriter, r *http.Request) error {
	grantType := r.URL.Query().Get("grant_type")
	if grantType == "" {
		// Try to read from body
		grantType = "password"
	}

	switch grantType {
	case "password":
		return a.tokenByPassword(w, r)
	case "refresh_token":
		return a.tokenByRefreshToken(w, r)
	default:
		return badRequestError("unsupported grant_type: %s", grantType)
	}
}

func (a *API) tokenByPassword(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	user, err := a.queries.GetUserByEmail(r.Context(), &req.Email)
	if err != nil {
		return unauthorizedError("invalid email or password")
	}

	if user.EncryptedPassword == nil {
		return unauthorizedError("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.EncryptedPassword), []byte(req.Password)); err != nil {
		return unauthorizedError("invalid email or password")
	}

	if !user.IsActive {
		return forbiddenError("account is deactivated")
	}

	// Update last sign in
	if err := a.queries.SetUserLastSignIn(r.Context(), user.ID); err != nil {
		log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("failed to update last sign in")
	}

	tokenResp, err := a.createSession(r, user.ID)
	if err != nil {
		return err
	}
	tokenResp.User = sanitizeUser(user)

	return sendJSON(w, http.StatusOK, tokenResp)
}

func (a *API) tokenByRefreshToken(w http.ResponseWriter, r *http.Request) error {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	// Look up refresh token
	token, err := a.queries.GetRefreshTokenByToken(r.Context(), req.RefreshToken)
	if err != nil {
		return unauthorizedError("invalid refresh token")
	}

	// If already revoked, revoke the entire family (token reuse detection)
	if token.Revoked {
		if token.SessionID != nil {
			if err := a.queries.RevokeRefreshTokenFamily(r.Context(), token.SessionID); err != nil {
				log.Warn().Err(err).Msg("failed to revoke refresh token family")
			}
		}
		return unauthorizedError("refresh token has been revoked")
	}

	// Revoke current token
	if err := a.queries.RevokeRefreshToken(r.Context(), req.RefreshToken); err != nil {
		log.Warn().Err(err).Msg("failed to revoke current refresh token")
	}

	// Issue new tokens
	user, err := a.queries.GetUserByID(r.Context(), token.UserID)
	if err != nil {
		return unauthorizedError("user not found")
	}

	// Generate new refresh token
	newRefreshToken := generateToken()
	_, err = a.queries.CreateRefreshToken(r.Context(), token.SessionID, token.UserID, newRefreshToken, &req.RefreshToken)
	if err != nil {
		return internalServerError("failed to create refresh token")
	}

	accessToken, err := a.generateAccessToken(user.ID)
	if err != nil {
		return internalServerError("failed to generate access token")
	}

	return sendJSON(w, http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "bearer",
		ExpiresIn:    a.config.JWT.Exp,
		User:         sanitizeUser(user),
	})
}

// Logout handles POST /auth/logout
func (a *API) Logout(w http.ResponseWriter, r *http.Request) error {
	userID := getUserID(r.Context())
	if err := a.queries.DeleteUserSessions(r.Context(), userID); err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("failed to delete user sessions")
	}
	return sendEmpty(w, http.StatusNoContent)
}

// GetUser handles GET /auth/user
func (a *API) GetUser(w http.ResponseWriter, r *http.Request) error {
	userID := getUserID(r.Context())
	user, err := a.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		return notFoundError("user not found")
	}
	return sendJSON(w, http.StatusOK, sanitizeUser(user))
}

// UpdateUser handles PUT /auth/user
func (a *API) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	userID := getUserID(r.Context())

	var body struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Display   *string `json:"display_name"`
		AvatarURL *string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return badRequestError("invalid request body")
	}

	user, err := a.queries.UpdateUser(r.Context(), userID, body.FirstName, body.LastName, body.Display, body.AvatarURL)
	if err != nil {
		return internalServerError("failed to update user")
	}
	return sendJSON(w, http.StatusOK, sanitizeUser(user))
}

// Recover handles POST /auth/recover
func (a *API) Recover(w http.ResponseWriter, r *http.Request) error {
	var req RecoverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	// Always return success to prevent email enumeration
	user, err := a.queries.GetUserByEmail(r.Context(), &req.Email)
	if err != nil {
		return sendJSON(w, http.StatusOK, map[string]string{"message": "if the email exists, a recovery link will be sent"})
	}

	token := generateToken()
	if err := a.queries.SetUserRecoveryToken(r.Context(), user.ID, token); err != nil {
		log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("failed to set recovery token")
	}

	// TODO: send recovery email via mailer

	return sendJSON(w, http.StatusOK, map[string]string{"message": "if the email exists, a recovery link will be sent"})
}

// Verify handles GET /auth/verify
func (a *API) Verify(w http.ResponseWriter, r *http.Request) error {
	tokenType := r.URL.Query().Get("type")
	token := r.URL.Query().Get("token")

	if token == "" {
		return badRequestError("missing token parameter")
	}

	switch tokenType {
	case "signup":
		user, err := a.queries.GetUserByConfirmationToken(r.Context(), token)
		if err != nil {
			return badRequestError("invalid or expired confirmation token")
		}
		if err := a.queries.SetUserEmailConfirmed(r.Context(), user.ID); err != nil {
			log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("failed to confirm user email")
		}
		if err := a.queries.SetUserConfirmationToken(r.Context(), user.ID, ""); err != nil {
			log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("failed to clear confirmation token")
		}

		// Redirect to frontend
		http.Redirect(w, r, a.config.SiteURL+"/auth/login?confirmed=true", http.StatusTemporaryRedirect)
		return nil

	case "recovery":
		user, err := a.queries.GetUserByRecoveryToken(r.Context(), token)
		if err != nil {
			return badRequestError("invalid or expired recovery token")
		}
		if err := a.queries.SetUserRecoveryToken(r.Context(), user.ID, ""); err != nil {
			log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("failed to clear recovery token")
		}

		// Redirect to frontend reset password page
		http.Redirect(w, r, a.config.SiteURL+"/auth/reset-password?token="+token, http.StatusTemporaryRedirect)
		return nil

	default:
		return badRequestError("unsupported verification type")
	}
}

// createSession creates a session, refresh token, and access token.
func (a *API) createSession(r *http.Request, userID uuid.UUID) (*TokenResponse, error) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		ip = r.RemoteAddr
	}
	ua := r.UserAgent()

	session, err := a.queries.CreateSession(r.Context(), userID, &ip, &ua)
	if err != nil {
		log.Error().Err(err).Str("ip", ip).Msg("CreateSession failed")
		return nil, internalServerError("failed to create session")
	}

	refreshToken := generateToken()
	_, err = a.queries.CreateRefreshToken(r.Context(), &session.ID, userID, refreshToken, nil)
	if err != nil {
		return nil, internalServerError("failed to create refresh token")
	}

	accessToken, err := a.generateAccessToken(userID)
	if err != nil {
		return nil, internalServerError("failed to generate access token")
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		ExpiresIn:    a.config.JWT.Exp,
	}, nil
}

// generateAccessToken creates a signed JWT.
func (a *API) generateAccessToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(a.config.JWT.Exp) * time.Second).Unix(),
		"iss": a.config.JWT.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.config.JWT.Secret))
}

func generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// sanitizeUser returns a safe user representation for API responses.
func sanitizeUser(user interface{}) map[string]interface{} {
	// We'll use JSON marshaling to convert the sqlc struct to a map,
	// then remove sensitive fields.
	data, _ := json.Marshal(user)
	var result map[string]interface{}
	_ = json.Unmarshal(data, &result)

	delete(result, "encrypted_password")
	delete(result, "confirmation_token")
	delete(result, "recovery_token")
	delete(result, "app_metadata")

	return result
}
