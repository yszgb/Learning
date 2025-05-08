package app

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"learning.com/learning/app/api"
	"learning.com/learning/models"
	"learning.com/learning/utils/helper"
	"learning.com/learning/utils/jsonlog"
)

type Application struct {
	config *models.Config
	helper *helper.Helper
	server *http.Server
	api    *api.Api
}

func New(configPath, currentDir string) (*Application, error) {
	config, err := models.NewConfig(configPath, currentDir)
	if err != nil {
		return nil, err
	}

	prettyLog := config.Env == "dev"

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo, prettyLog)
	helper := helper.New(logger)
	api, err := api.New(helper, config)
	if err != nil {
		return nil, err
	}

	return &Application{
		config: config,
		helper: helper,
		api:    api,
	}, nil
}

// 启动非视频服务
func (app *Application) runServer() error {
	tlsCfg := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	addr := fmt.Sprintf("%s:%d", app.config.Server, app.config.Port)
	app.server = &http.Server{
		Addr:      addr,
		Handler:   app.api.Routes(),
		TLSConfig: tlsCfg,

		IdleTimeout:    time.Minute,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 8192,
	}

	app.helper.Logger.Info("Starting server", map[string]string{
		"addr": addr,
		"env":  app.config.Env,
	})

	return app.server.ListenAndServeTLS(app.config.Path.CertPath(), app.config.Path.KeyPath())
}

func (app *Application) Run() error {
	return app.runServer()
}
