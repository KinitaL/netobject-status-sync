package interfaces

import (
	"netobject-status-sync/internal/config"
	"netobject-status-sync/internal/model"
)

type Store interface {
	ConnectToDb(config *config.Config) error
	FindMacs() ([]model.Device, error)
	SetInstalledStatus(device model.Device, status bool) error
}
