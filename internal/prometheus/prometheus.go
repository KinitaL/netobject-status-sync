package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"netobject-status-sync/internal/config"
	"netobject-status-sync/internal/model"
)

type Prometheus struct {
	Addr        string
	MacsToSend  prometheus.Counter
	HandledMacs prometheus.Counter
}

func NewPrometheus(config *config.Config) *Prometheus {
	return &Prometheus{
		MacsToSend: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "macs_to_send",
				Help: "Количество маков, отправленных в канал",
			}),
		HandledMacs: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "installed_macs",
				Help: "Количество маков, обработанных через papi, отправленных в канал",
			}),
		Addr: config.Prometheus.Addr,
	}
}

func (prom *Prometheus) AddHandledMacs(macsWithStatuses map[model.Device]bool) {
	prom.HandledMacs.Add(float64(len(macsWithStatuses)))
}

func (prom *Prometheus) AddMacsToSend(chunk []model.Device) {
	prom.MacsToSend.Add(float64(len(chunk)))
}

func (prom *Prometheus) Register() {
	prometheus.MustRegister(prom.MacsToSend)
	prometheus.MustRegister(prom.HandledMacs)
}

func (prom *Prometheus) Listen() {
	defer func() {
		if err := recover(); err != nil {
			log.Warnf("error prometheus listening: %s\n ", err)
			go prom.Listen()
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(prom.Addr, nil); err != nil {
		log.Warnf("error prometheus listening: %s\n ", err)
	}
}
