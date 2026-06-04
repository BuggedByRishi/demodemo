package handler

import (
	"net/http"
	"strconv"

	"Tally-server/internal/model"
	"Tally-server/internal/repository"
)

type ReferenceHandler struct{ repo repository.ReferenceRepository }

func NewReferenceHandler(repo repository.ReferenceRepository) *ReferenceHandler {
	return &ReferenceHandler{repo: repo}
}

// GET /countries
func (h *ReferenceHandler) ListCountries(w http.ResponseWriter, r *http.Request) {
	countries, err := h.repo.ListCountries(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if countries == nil {
		countries = []model.Country{}
	}
	writeJSON(w, http.StatusOK, countries)
}

// GET /countries/{id}/states
func (h *ReferenceHandler) ListStates(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path) // countries/{id}/states
	if len(parts) < 2 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}
	countryID, err := strconv.Atoi(parts[len(parts)-2])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid country_id")
		return
	}
	states, err := h.repo.ListStatesByCountry(r.Context(), countryID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if states == nil {
		states = []model.State{}
	}
	writeJSON(w, http.StatusOK, states)
}

// GET /themes
func (h *ReferenceHandler) ListThemes(w http.ResponseWriter, r *http.Request) {
	themes, err := h.repo.ListThemes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if themes == nil {
		themes = []model.FrontendTheme{}
	}
	writeJSON(w, http.StatusOK, themes)
}

