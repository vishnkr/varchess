package game

import (
	"net/http"
	"varchess/internal/user"

	"github.com/go-chi/chi/v5"
)
type API struct{
	gameService Service
	userService user.Service
}

func NewAPI(gameService Service,userService user.Service)API{
	return API{gameService, userService}
}

// endpoint POST /templates
func (api *API) HandleCreateGame(w http.ResponseWriter, r *http.Request){
	//ctx := r.Context()
	//logger := logger.FromContext(ctx)

}

func (api *API) HandleGetGames(w http.ResponseWriter, r *http.Request){
	//ctx := r.Context()
	//logger := logger.FromContext(ctx)
}

func (api *API) HandleGetGame(w http.ResponseWriter, r *http.Request){
	//ctx := r.Context()
	//logger := logger.FromContext(ctx)
}

func (api *API) RegisterHandlers(r chi.Router){
	r.Route("/games",func(r chi.Router){
		r.Group(func(r chi.Router){
			r.Get("/",api.HandleGetGames)
			r.Route("/{id}",func(r chi.Router){
				r.Get("/",api.HandleGetGame)
			})
			r.Post("/save",api.HandleCreateGame)
		})
	})
}