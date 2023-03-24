package internal

import (
	"github.com/google/uuid"
	"k8s.io/klog/v2"
	"ms-keys/pkg"
)

type LogTransport struct {
}

func (l *LogTransport) Send(register pkg.RegisterData, session uuid.UUID) {
	host := "http://localhost:8899"
	url := host + "/api/verify?session=" + session.String()
	klog.Info("Register verify url: ", url)
}
