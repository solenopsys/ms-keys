package main

import "github.com/google/uuid"

type MailServ struct {
	host string
	port string
}

func (m *MailServ) sendMail(register RegisterData, session uuid.UUID) {
	// send email

}
