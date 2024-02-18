package template

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)
type API struct{
	service Service
}

func NewAPI(service Service)API{
	return API{service}
}

const (
	errInternalServer = "Internal Server Error"
	errTemplateNotFound = "Template not found"
	errInvalidTemplateID = "Invalid TemplateId"
)
// endpoint POST /templates
func (api *API) HandleCreateTemplate(w http.ResponseWriter, r *http.Request){
	//ctx := r.Context()
	//logger := logger.FromContext(ctx)

}

func (api *API) HandleGetTemplates(w http.ResponseWriter, r *http.Request){
	//ctx := r.Context()
	//logger := logger.FromContext(ctx)
	/*
	userID:= middleware.getUserIDFromContext(ctx)
	templates,err := api.service.GetTemplatesByUser(ctx,userID)
	if err!=nil{
		return
	}
	Write response with status
	*/
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