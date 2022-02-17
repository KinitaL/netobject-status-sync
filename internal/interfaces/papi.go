package interfaces

import (
	"netobject-status-sync/internal/model"
)

type Papi interface {
	Send(devices []model.Device, result chan map[model.Device]bool, prom Prometheus)
	SetStatuses(devices []model.Device, installedMacs []string) map[model.Device]bool
}
