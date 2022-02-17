package application

import (
	"netobject-status-sync/internal/config"
	"netobject-status-sync/internal/papi"
	"netobject-status-sync/internal/prometheus"
	"netobject-status-sync/internal/store"
	"netobject-status-sync/internal/syncer"
)

type Application struct {
	config     *config.Config
	store      *store.Store
	papi       *papi.Papi
	syncer     *syncer.Syncer
	prometheus *prometheus.Prometheus
}

func NewApp() *Application {
	return &Application{
		store: store.NewStore(),
	}
}

func (app *Application) Configure() error {
	cfg, err := config.ReadConfig("app.yaml")
	if err != nil {
		return err
	}
	app.config = cfg

	if err = app.store.ConnectToDb(app.config); err != nil {
		return err
	}

	app.papi = papi.NewPapi(app.config, "")

	app.prometheus = prometheus.NewPrometheus(app.config)
	app.prometheus.Register()

	app.syncer = syncer.NewSyncer(app.papi, app.store, app.config, app.prometheus)

	return nil
}

func (app *Application) Run() {
	go app.prometheus.Listen()
	app.syncer.Run(app.config.App.StartHour, app.config.App.StartMinute, app.config.App.StartSecond, app.config.App.SleepAfterSync)
}
