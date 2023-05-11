package internal

import (
	"github.com/google/uuid"
	klog "k8s.io/klog/v2"
	"ms-keys/pkg"
	"net/smtp"
	"strings"
)

type MailTransport struct {
	From     string
	Host     string
	Port     string
	AuthHost string
	Password string
	Username string
}

func (m *MailTransport) Send(register pkg.RegisterData, session uuid.UUID) {
	to := []string{register.Login}
	subject := "Register in Solenopsys"
	body := "You are successfully registered in Solenopsys. Please click on the link below to verify your email address.\r\n" +
		m.AuthHost + "/verify?session=" + session.String()

	// Set up authentication information.
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	// Set up the email message.
	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body

	// Send the email message.
	err := smtp.SendMail(m.Host+":"+m.Port, auth, m.From, to, []byte(msg))
	if err != nil {
		klog.Error(err)
	}
}
