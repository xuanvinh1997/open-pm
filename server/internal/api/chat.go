package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/service"
	"github.com/rs/zerolog/log"
)

type SendChatRequest struct {
	Message string `json:"message" validate:"required,min=1,max=2000"`
}

// SendChatMessage handles POST .../projects/{projectID}/chat
func (a *API) SendChatMessage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	if a.llm == nil {
		return badRequestError("AI chat is not enabled")
	}

	var req SendChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	// Store user message
	_, err := a.queries.CreateChatMessage(ctx, projectID, workspaceID, userID, "user", req.Message)
	if err != nil {
		log.Error().Err(err).Msg("failed to store user chat message")
		return internalServerError("failed to process message")
	}

	// Load recent chat history for context
	history, err := a.queries.ListChatMessages(ctx, projectID, userID, 20)
	if err != nil {
		log.Error().Err(err).Msg("failed to load chat history")
		history = nil
	}

	// Build project context for the system prompt
	systemPrompt, err := a.buildProjectContext(ctx, projectID)
	if err != nil {
		log.Error().Err(err).Msg("failed to build project context")
		systemPrompt = "You are an AI assistant for a project management tool. Answer questions about the project based on the conversation."
	}

	// Build LLM messages from history
	var llmMessages []service.LLMMessage
	for _, msg := range history {
		llmMessages = append(llmMessages, service.LLMMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Call LLM
	response, err := a.llm.Chat(ctx, systemPrompt, llmMessages)
	if err != nil {
		log.Error().Err(err).Msg("LLM chat failed")
		return internalServerError("failed to get AI response")
	}

	// Store assistant response
	assistantMsg, err := a.queries.CreateChatMessage(ctx, projectID, workspaceID, userID, "assistant", response)
	if err != nil {
		log.Error().Err(err).Msg("failed to store assistant chat message")
		return internalServerError("failed to save response")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"message": assistantMsg,
	})
}

// ListChatHistory handles GET .../projects/{projectID}/chat
func (a *API) ListChatHistory(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	messages, err := a.queries.ListChatMessages(ctx, projectID, userID, limit)
	if err != nil {
		return internalServerError("failed to load chat history")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"results": messages,
	})
}

// ClearChatHistory handles DELETE .../projects/{projectID}/chat
func (a *API) ClearChatHistory(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)

	if err := a.queries.DeleteChatHistory(ctx, projectID, userID); err != nil {
		return internalServerError("failed to clear chat history")
	}

	return sendEmpty(w, http.StatusNoContent)
}

// buildProjectContext creates a system prompt with project data for LLM context.
func (a *API) buildProjectContext(ctx context.Context, projectID uuid.UUID) (string, error) {
	project, err := a.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return "", err
	}

	// Get issue stats
	byState, _ := a.queries.CountIssuesByState(ctx, projectID)
	byPriority, _ := a.queries.CountIssuesByPriority(ctx, projectID)
	overdueIssues, _ := a.queries.ListOverdueIssues(ctx, projectID)
	blockedIDs, _ := a.queries.ListBlockedIssueIDs(ctx, projectID)

	// Get members
	members, _ := a.queries.ListProjectMembers(ctx, projectID)

	// Get cycles
	cycles, _ := a.queries.ListCyclesByProject(ctx, projectID)

	// Build context string
	ctxStr := "You are an AI assistant for the project management tool Open-PM. " +
		"Answer questions about the project based on the following data. " +
		"Be concise and helpful. Use the data provided to give accurate answers.\n\n"

	ctxStr += "## Project: " + project.Name + "\n"
	if project.Description != "" {
		ctxStr += "Description: " + project.Description + "\n"
	}
	ctxStr += "\n"

	// Members
	if len(members) > 0 {
		ctxStr += "## Team Members\n"
		for _, m := range members {
			name := m.DisplayName
			if name == "" {
				name = m.FirstName + " " + m.LastName
			}
			ctxStr += "- " + name + "\n"
		}
		ctxStr += "\n"
	}

	// Issue stats by state
	if len(byState) > 0 {
		ctxStr += "## Issues by State\n"
		var total int
		for _, s := range byState {
			ctxStr += "- " + s.StateName + " (" + s.StateGroup + "): " + strconv.Itoa(s.Count) + "\n"
			total += s.Count
		}
		ctxStr += "- Total: " + strconv.Itoa(total) + "\n\n"
	}

	// Issue stats by priority
	if len(byPriority) > 0 {
		ctxStr += "## Issues by Priority\n"
		for _, p := range byPriority {
			ctxStr += "- " + p.Priority + ": " + strconv.Itoa(p.Count) + "\n"
		}
		ctxStr += "\n"
	}

	// Overdue issues
	if len(overdueIssues) > 0 {
		ctxStr += "## Overdue Issues (" + strconv.Itoa(len(overdueIssues)) + ")\n"
		for _, i := range overdueIssues {
			ctxStr += "- " + i.Name
			if i.TargetDate != nil {
				ctxStr += " (due: " + i.TargetDate.Format("2006-01-02") + ")"
			}
			ctxStr += "\n"
		}
		ctxStr += "\n"
	}

	// Blocked issues
	if len(blockedIDs) > 0 {
		ctxStr += "## Blocked Issues: " + strconv.Itoa(len(blockedIDs)) + "\n\n"
	}

	// Cycles
	if len(cycles) > 0 {
		ctxStr += "## Cycles/Sprints\n"
		for _, c := range cycles {
			ctxStr += "- " + c.Name
			if c.StartDate != nil && c.EndDate != nil {
				ctxStr += " (" + c.StartDate.Format("2006-01-02") + " to " + c.EndDate.Format("2006-01-02") + ")"
			}
			ctxStr += "\n"
		}
		ctxStr += "\n"
	}

	return ctxStr, nil
}
