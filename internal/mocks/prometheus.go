package mocks

import "netobject-status-sync/internal/model"

type MockProm struct{}

func (prom *MockProm) Register()                                             {}
func (prom *MockProm) Listen()                                               {}
func (prom *MockProm) AddHandledMacs(macsWithStatuses map[model.Device]bool) {}
func (prom *MockProm) AddMacsToSend(chunk []model.Device)                    {}
