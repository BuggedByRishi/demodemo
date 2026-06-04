package main

import (
	"log"
	"net/http"

	"Tally-server/internal/config"
	"Tally-server/internal/db"
	"Tally-server/internal/handler"
	"Tally-server/internal/repository"
	"Tally-server/internal/service"
)

func main() {
	// ── 1. Load config from DB_* env vars ─────────────────────────────────
	cfg := config.Load()

	// ── 2. Connect to PostgreSQL ───────────────────────────────────────────
	database, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer database.Close()
	log.Println("database connected successfully")

	// ── 3. Wire layers: repository → service → handler ────────────────────
	orgRepo := repository.NewOrganizationRepository(database)
	agentRepo := repository.NewAgentRepository(database)
	refRepo := repository.NewReferenceRepository(database)

	orgSvc := service.NewOrganizationService(orgRepo)
	agentSvc := service.NewAgentService(agentRepo)

	healthH := handler.NewHealthHandler(database)
	orgH := handler.NewOrganizationHandler(orgSvc)
	agentH := handler.NewAgentHandler(agentSvc)
	refH := handler.NewReferenceHandler(refRepo)

	// ── 4. Register routes ─────────────────────────────────────────────────
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", healthH.Health)

	// Reference data (read-only, seeded by migrations)
	mux.HandleFunc("GET /countries", refH.ListCountries)
	mux.HandleFunc("GET /countries/{id}/states", refH.ListStates)
	mux.HandleFunc("GET /themes", refH.ListThemes)

	// Organizations
	mux.HandleFunc("POST /organizations", orgH.Create)
	mux.HandleFunc("GET /organizations", orgH.List)
	mux.HandleFunc("GET /organizations/{id}", orgH.GetByID)

	// Agents
	mux.HandleFunc("POST /agents", agentH.Create)
	mux.HandleFunc("GET /agents/{id}", agentH.GetByID)
	mux.HandleFunc("GET /organizations/{id}/agents", agentH.ListByOrganization)

	// ── 5. Start server ────────────────────────────────────────────────────
	addr := ":" + cfg.Port
	log.Printf("server starting on %s (env: %s)", addr, cfg.Env)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
