package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"Tally-server/internal/model"
	"Tally-server/internal/service"
)

type OrganizationHandler struct{ svc *service.OrganizationService }

func NewOrganizationHandler(svc *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{svc: svc}
}

// POST /organizations
func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	org, err := h.svc.Create(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, org)
}

// GET /organizations
func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if orgs == nil {
		orgs = []model.Organization{}
	}
	writeJSON(w, http.StatusOK, orgs)
}

// GET /organizations/{id}
func (h *OrganizationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := lastPathSegment(r)
	org, err := h.svc.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		writeError(w, http.StatusNotFound, "organization not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, org)
}

