package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"learning.com/learning/app/middleware"
	"learning.com/learning/models"
	"learning.com/learning/utils/helper"
)

type Api struct {
	config     *models.Config
	helper     *helper.Helper
	models     *models.Models
	middleware *middleware.Middleware
}

func New(helper *helper.Helper, config *models.Config) (*Api, error) {
	models, err := models.NewModels(&config.Database, helper)
	if err != nil {
		return nil, err
	}
	middleware := middleware.New(helper, &config.Limiter, models)
	return &Api{
		config:     config,
		helper:     helper,
		models:     models,
		middleware: middleware,
	}, nil
}

func (api *Api) Version() string {
	return "v1"
}

func (api *Api) Routes() http.Handler {
	// 初始化路由
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(api.helper.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(api.helper.MethodNotAllowedResponse)

	api.utilsRoutes(router)
	api.userRoutes(router)

	standard := alice.New(
		api.middleware.Metrics,
		api.middleware.RecoverPanic,
		api.middleware.LogRequest,
		api.middleware.RateLimit,
		api.middleware.Authenticate,
	)

	return standard.Then(router)
}
