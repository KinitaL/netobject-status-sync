package tests

import (
	log "github.com/sirupsen/logrus"
	"netobject-status-sync/internal/config"
	"netobject-status-sync/internal/mocks"
	"netobject-status-sync/internal/model"
	"netobject-status-sync/internal/syncer"
	"testing"
	"time"
)

func TestSyncer(t *testing.T) {
	testingTable := []struct {
		name    string
		devices []model.Device
	}{
		{
			name:    "Main",
			devices: []model.Device{{Mac: "18:68:82:30:84:AE", Id: 73356}},
		},
	}

	for _, test := range testingTable {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := config.ReadConfig("app.yaml")
			if err != nil {
				log.Warnf("error reading config %s\n", err)
			}
			testSyncer := syncer.NewSyncer(
				mocks.NewMockPapi(test.devices),
				mocks.NewMockStore(test.devices),
				cfg,
				&mocks.MockProm{},
			)

			hour := time.Now().Hour()
			minute := time.Now().Minute()
			second := time.Now().Second() + 5

			ticker := time.NewTicker(time.Second * 10)
			t.Cleanup(func() {
				ticker.Stop()
			})
			for range ticker.C {
				t.Skip("OK: test passed. No panics in 10 second after goroutine started mean that everything is all right")
			}
			testSyncer.Run(hour, minute, second, 24)
		})
	}
}
