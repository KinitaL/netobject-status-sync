package mocks

import (
	"netobject-status-sync/internal/config"
	"netobject-status-sync/internal/model"
)

type MockStore struct {
	devices []model.Device
}

func NewMockStore(devices []model.Device) *MockStore {
	return &MockStore{devices: devices}
}

func (mock *MockStore) ConnectToDb(config *config.Config) error {
	return nil
}

func (mock *MockStore) FindMacs() ([]model.Device, error) {
	return mock.devices, nil
}

func (mock *MockStore) SetInstalledStatus(device model.Device, status bool) error {
	return nil
}
