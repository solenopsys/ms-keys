package internal

import (
	"github.com/google/uuid"
	"ms-keys/pkg"
)

type MailSenderInterface interface {
	SendEMail(register pkg.RegisterData, session uuid.UUID)
}

type PersistedData interface {
	LoadRegister(email string) (pkg.RegisterData, error)
	SaveRegister(register pkg.RegisterData)
}
