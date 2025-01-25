package api

import (
	"vc-core/internal/api/handlers"
	"vc-core/internal/db"
	"vc-core/internal/logger"

	"github.com/go-chi/chi/v5"
)

type API struct {
	logger *logger.Logger
	db     *db.DB
}

func NewAPI(logger *logger.Logger, db *db.DB) *API {
	return &API{logger: logger, db: db}
}

func (api *API) RegisterHandlers(r chi.Router, db *db.DB) {
	r.Route("/games", func(r chi.Router) {
		r.Get("/", handlers.HandleGetGames(db))
		r.Post("/", handlers.HandleCreateGame(db))
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.HandleGetGame(db))
		})
	})
	r.Route("/templates", func(r chi.Router) {
		r.Get("/", handlers.HandleGetTemplates(db))
		r.Post("/", handlers.HandleCreateTemplate(db))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.HandleGetTemplate(db))
			r.Put("/", handlers.HandleUpdateTemplate(db))
			r.Delete("/", handlers.HandleDeleteTemplate(db))
		})
	})
	r.Route("/settings", func(r chi.Router) {
		r.Get("/", handlers.HandleGetSettings(db))
		r.Put("/", handlers.HandleUpdateSettings(db))
	})

	/*r.Route("/ws", func(r chi.Router) {
		r.Get("/", ws.HandleWebSocket)
	})*/
}
