package papi

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"netobject-status-sync/internal/config"
	"netobject-status-sync/internal/interfaces"
	"netobject-status-sync/internal/model"
	"time"
)

type Papi struct {
	BaseUrl          string
	Url              string
	client           *resty.Client
	SecondSleep      int
	NumberOfAttempts int
}

func NewPapi(config *config.Config, url string) *Papi {
	if url == "" {
		url = config.Papi.Url
	}
	return &Papi{
		BaseUrl:          config.Papi.BaseUrl,
		Url:              url,
		client:           resty.New(),
		SecondSleep:      config.App.SecondSleep,
		NumberOfAttempts: config.App.NumberOfAttempts,
	}
}

func (papi *Papi) Send(devices []model.Device, result chan map[model.Device]bool, prom interfaces.Prometheus) {
	defer func() {
		if err := recover(); err != nil {
			log.Warnf("Something went wrong sending to papi: %s\n", err)
			return
		}
	}()

	var macs []string
	for _, device := range devices {
		macs = append(macs, device.Mac)
	}

	resp, err := papi.client.R().SetBody(macs).Post(papi.BaseUrl + papi.Url)
	if err != nil {
		log.Warn(err)
		ticker := time.NewTicker(time.Second * time.Duration(papi.SecondSleep))
		defer ticker.Stop()
	reSendLoop:
		for i := 0; i < papi.NumberOfAttempts; i++ {
			for range ticker.C {
				resp, err = papi.client.R().SetBody(macs).Post(papi.BaseUrl + papi.Url)
				if err != nil {
					continue reSendLoop
				}
				break reSendLoop
			}
		}
	}
	var installedMacs []string
	if err = json.Unmarshal(resp.Body(), &installedMacs); err != nil {
		log.Warnf("Something went wrong encoding response from papi: %s\n", err)
		return
	}

	macsWithStatuses := papi.SetStatuses(devices, installedMacs)
	result <- macsWithStatuses
	prom.AddHandledMacs(macsWithStatuses)
	return
}

func (papi *Papi) SetStatuses(devices []model.Device, installedMacs []string) map[model.Device]bool {
	macsWithStatuses := make(map[model.Device]bool)
	for _, device := range devices {
		status := false
		for _, installedMac := range installedMacs {
			if installedMac == device.Mac {
				status = true
			}
		}
		macsWithStatuses[device] = status
	}
	return macsWithStatuses
}
