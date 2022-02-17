package mocks

import (
	"netobject-status-sync/internal/interfaces"
	"netobject-status-sync/internal/model"
)

type MockPapi struct {
	expectedDevicesMacs []string
}

func NewMockPapi(devices []model.Device) *MockPapi {
	var expectedDevicesMacs []string
	for _, device := range devices {
		expectedDevicesMacs = append(expectedDevicesMacs, device.Mac)
	}
	return &MockPapi{expectedDevicesMacs: expectedDevicesMacs}
}

func (mock *MockPapi) Send(devices []model.Device, result chan map[model.Device]bool, prom interfaces.Prometheus) {
	var installedMacs []string
	for _, device := range devices {
		installedMacs = append(installedMacs, device.Mac)
	}
	macsWithStatuses := mock.SetStatuses(devices, installedMacs)
	result <- macsWithStatuses
}

func (mock *MockPapi) SetStatuses(devices []model.Device, installedMacs []string) map[model.Device]bool {
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

	//тестовая проверка, которая вызывает панику в случае, если приходить некорректный результат.
	for _, device := range devices {
		for _, expectedMac := range mock.expectedDevicesMacs {
			if expectedMac == device.Mac {
				continue
			} else {
				panic("FAILED: Unexpected device")
			}
		}
	}

	return macsWithStatuses
}
