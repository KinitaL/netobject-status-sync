package syncer

import (
	log "github.com/sirupsen/logrus"
	"netobject-status-sync/internal/config"
	"netobject-status-sync/internal/interfaces"
	"netobject-status-sync/internal/model"
	"sync"
	"time"
)

type Syncer struct {
	papi                  interfaces.Papi
	store                 interfaces.Store
	chunksCh              chan []model.Device
	devicesWithStatusesCh chan map[model.Device]bool
	maxItemsToSend        int
	prom                  interfaces.Prometheus
}

func NewSyncer(papi interfaces.Papi, store interfaces.Store, config *config.Config, prom interfaces.Prometheus) *Syncer {
	return &Syncer{
		papi:                  papi,
		store:                 store,
		chunksCh:              make(chan []model.Device, 10000),
		devicesWithStatusesCh: make(chan map[model.Device]bool, 20000),
		maxItemsToSend:        config.App.MaxItemsToSend,
		prom:                  prom,
	}
}

func (syncer *Syncer) Run(hour, min, sec, afterSync int) {
	var wg sync.WaitGroup
	wg.Add(1)

	go syncer.GetMacs(hour, min, sec, afterSync)
	go syncer.Send()
	go syncer.SetStatuses()

	wg.Wait()
}

func (syncer *Syncer) GetMacs(hour, min, sec, afterSync int) {
	defer func() {
		//Запускам горутину снова, так как сервис должен работать непрерывно.
		if err := recover(); err != nil {
			log.Warnf("function GetMacs ended up with error: %s\n", err)
			log.Warn("restart GetMacs")
			go syncer.GetMacs(hour, min, sec, afterSync)
		}
	}()

	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Warnf("Can't load location: %s\n", err)
	}

	now := time.Now().Local()
	firstCallTime := time.Date(
		now.Year(), now.Month(), now.Day(), hour, min, sec, 0, loc)
	if firstCallTime.Before(now) {
		firstCallTime = firstCallTime.Add(time.Hour * time.Duration(afterSync))
	}

	duration := firstCallTime.Sub(time.Now().Local())

	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		for range ticker.C {
			ticker.Reset(time.Hour * time.Duration(afterSync))
			devices, err := syncer.store.FindMacs()
			if err != nil {
				log.Warnf("error find macs from database: %s\n", err)
			}

			var chunk []model.Device
			for _, device := range devices {
				chunk = append(chunk, device)
				if len(chunk) > syncer.maxItemsToSend {
					syncer.chunksCh <- chunk
					syncer.prom.AddMacsToSend(chunk)
					chunk = nil
				}
			}

			if chunk != nil {
				syncer.chunksCh <- chunk
				syncer.prom.AddMacsToSend(chunk)
				chunk = nil
			}
		}
	}
}

func (syncer *Syncer) Send() {
	defer func() {
		if err := recover(); err != nil {
			log.Warnf("function Send ended up with error: %s\n", err)
			log.Warn("restart Send")
			go syncer.Send()
		}
	}()

	for {
		for macsChunk := range syncer.chunksCh {
			go syncer.papi.Send(macsChunk, syncer.devicesWithStatusesCh, syncer.prom)
		}
	}
}

func (syncer *Syncer) SetStatuses() {
	defer func() {
		if err := recover(); err != nil {
			log.Warnf("function SetStatuses ended up with error: %s\n", err)
			log.Warn("restart SetStatuses")
			go syncer.SetStatuses()
		}
	}()

	for {
		for deviceWithStatuses := range syncer.devicesWithStatusesCh {
			// избежать гонок за чтение
			deviceWithStatuses := deviceWithStatuses
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Warnf("error set statuses to database: %s\n", err)
						return
					}
				}()

				for device, status := range deviceWithStatuses {
					if err := syncer.store.SetInstalledStatus(device, status); err != nil {
						log.Warnf("error set statuses to database: %s\n", err)
					}
				}

				return
			}()
		}
	}
}
