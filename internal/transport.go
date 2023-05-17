package internal

import (
	"github.com/google/uuid"
	"ms-keys/pkg"
)

type Transport interface {
	Send(register pkg.RegisterData, session uuid.UUID)
}

type Service struct {
	transports map[string]Transport
}

func (t *Service) AddTransport(name string, transport Transport) {
	t.transports[name] = transport
}

func (t *Service) Send(name string, register pkg.RegisterData, session uuid.UUID) {
	go t.transports[name].Send(register, session)
}

func (t *Service) List() []string {
	keys := make([]string, 0, len(t.transports))
	for k := range t.transports {
		keys = append(keys, k)
	}
	return keys
}

func NewTransportHolder() *Service {
	th := &Service{
		transports: make(map[string]Transport),
	}
	return th
}
