package template

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
func (api *API) HandleCreateTemplate(w http.ResponseWriter, r *http.Request){
	//ctx := r.Context()
	//logger := logger.FromContext(ctx)

}

func (api *API) HandleGetTemplates(w http.ResponseWriter, r *http.Request){
	//ctx := r.Context()
	//logger := logger.FromContext(ctx)
}

func (api *API) HandleGetTemplate(w http.ResponseWriter, r *http.Request){
	//ctx := r.Context()
	//logger := logger.FromContext(ctx)
}

func (api *API) RegisterHandlers(r chi.Router){
	r.Route("/templates",func(r chi.Router){
		r.Group(func(r chi.Router){
			r.Get("/",api.HandleGetTemplates)
			r.Route("/{id}",func(r chi.Router){
				r.Get("/",api.HandleGetTemplate)
			})
			r.Post("/",api.HandleCreateTemplate)
		})
	})
}