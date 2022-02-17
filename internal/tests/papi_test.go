package tests

import (
	log "github.com/sirupsen/logrus"
	"netobject-status-sync/internal/config"
	"netobject-status-sync/internal/model"
	"netobject-status-sync/internal/papi"
	"reflect"
	"testing"
)

func TestSetStatuses(t *testing.T) {
	testingTable := []struct {
		name           string
		devices        []model.Device
		installedMacs  []string
		expectedResult map[model.Device]bool
	}{
		{
			name: "All installed",
			devices: []model.Device{
				{Mac: "18:68:82:30:84:AE", Id: 73356},
				{Mac: "18:68:82:30:2C:C1", Id: 17583},
			},
			installedMacs: []string{
				"18:68:82:30:84:AE",
				"18:68:82:30:2C:C1",
			},
			expectedResult: map[model.Device]bool{
				{Mac: "18:68:82:30:84:AE", Id: 73356}: true,
				{Mac: "18:68:82:30:2C:C1", Id: 17583}: true,
			},
		},
		{
			name: "All devices wasn't installed by papi",
			devices: []model.Device{
				{Mac: "18:68:82:30:84:AE", Id: 73356},
				{Mac: "18:68:82:30:2C:C1", Id: 17583},
			},
			installedMacs: []string{},
			expectedResult: map[model.Device]bool{
				{Mac: "18:68:82:30:84:AE", Id: 73356}: false,
				{Mac: "18:68:82:30:2C:C1", Id: 17583}: false,
			},
		},
		{
			name: "Only half of devices are installed",
			devices: []model.Device{
				{Mac: "18:68:82:30:84:AE", Id: 73356},
				{Mac: "18:68:82:30:2C:C1", Id: 17583},
			},
			installedMacs: []string{
				"18:68:82:30:84:AE",
			},
			expectedResult: map[model.Device]bool{
				{Mac: "18:68:82:30:84:AE", Id: 73356}: true,
				{Mac: "18:68:82:30:2C:C1", Id: 17583}: false,
			},
		},
	}

	cfg, err := config.ReadConfig("app.yaml")
	if err != nil {
		log.Warnf("error reading config %s\n", err)
	}
	testPapi := papi.NewPapi(cfg, "")

	for _, test := range testingTable {
		t.Run(test.name, func(t *testing.T) {
			result := testPapi.SetStatuses(test.devices, test.installedMacs)
			if !reflect.DeepEqual(result, test.expectedResult) {
				t.Error("Result of function isn't same with expected result")
			}
		})
	}
}
