package interfaces

import "netobject-status-sync/internal/model"

type Prometheus interface {
	AddHandledMacs(macsWithStatuses map[model.Device]bool)
	AddMacsToSend(chunk []model.Device)
	Register()
	Listen()
}
