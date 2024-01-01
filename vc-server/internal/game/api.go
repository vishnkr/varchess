package game

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)
type API struct{
	repository Repository
}

func NewAPI(repository Repository)API{
	return API{repository}
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
			r.Post("/",api.HandleCreateGame)
		})
	})
}