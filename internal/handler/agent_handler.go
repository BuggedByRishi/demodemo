package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"Tally-server/internal/model"
	"Tally-server/internal/service"
)

type AgentHandler struct{ svc *service.AgentService }

func NewAgentHandler(svc *service.AgentService) *AgentHandler {
	return &AgentHandler{svc: svc}
}

// POST /agents
func (h *AgentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateAgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	agent, err := h.svc.Create(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, agent)
}

// GET /agents/{id}
func (h *AgentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := lastPathSegment(r)
	agent, err := h.svc.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		writeError(w, http.StatusNotFound, "agent not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, agent)
}

// GET /organizations/{id}/agents
func (h *AgentHandler) ListByOrganization(w http.ResponseWriter, r *http.Request) {
	// path: /organizations/{org_id}/agents — extract second-to-last segment
	parts := splitPath(r.URL.Path)
	if len(parts) < 3 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}
	orgID := parts[len(parts)-2] // organizations/{org_id}/agents

	agents, err := h.svc.ListByOrganization(r.Context(), orgID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if agents == nil {
		agents = []model.Agent{}
	}
	writeJSON(w, http.StatusOK, agents)
}

func splitPath(path string) []string {
	var parts []string
	for _, p := range splitSlash(path) {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func splitSlash(s string) []string {
	var out []string
	cur := ""
	for _, c := range s {
		if c == '/' {
			out = append(out, cur)
			cur = ""
		} else {
			cur += string(c)
		}
	}
	out = append(out, cur)
	return out
}

