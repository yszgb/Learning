package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (api *Api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"status": "available",
		"system_info": map[string]string{
			"environment": api.config.Env,
			"version":     api.Version(),
		},
	}

	resp := api.helper.NewResponse(0, data)
	err := api.helper.WriteJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		api.helper.ServerErrorResponse(w, r, err)
	}
}

func (api *Api) utilsRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", api.healthCheckHandler)
}
