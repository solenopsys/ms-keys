package internal

import (
	"github.com/google/uuid"
	"k8s.io/klog/v2"
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
	klog.Info("Smtp: ", m.Host+":"+m.Port)
	klog.Info("User:", m.Username)
	klog.Info("Start send to:", to)
	err := smtp.SendMail(m.Host+":"+m.Port, auth, m.From, to, []byte(msg))
	klog.Info("End send to:", to)
	if err != nil {
		klog.Error(err)
	}
}

/*
func (m *MailTransport) Send(register pkg.RegisterData, session uuid.UUID) {

	qrCode, err := qrcode.New("Hello, World!", qrcode.Medium)
	if err != nil {
		panic(err)
	}
	// Generate QR code image.

	png, err := qrCode.PNG(100)
	qrCodeBase64 := base64.StdEncoding.EncodeToString(png)

	// Define email body.
	html := fmt.Sprintf(`<html>
	<body>
		<h1>Hello!</h1>
		<p>Here's a QR code:</p>
		<img src="cid:qrcode" alt="QR Code" />
	</body>
</html>`)
	htmlPart := make(map[string][]string)
	htmlPart["Content-Type"] = []string{"text/html; charset=UTF-8"}
	htmlPart["Content-Disposition"] = []string{"inline"}
	htmlPart["Content-ID"] = []string{"<qrcode>"}
	htmlPartBody := strings.NewReader(html)

	// Define image attachment.
	imagePart := make(map[string][]string)
	imagePart["Content-Type"] = []string{"image/png"}
	imagePart["Content-Disposition"] = []string{"attachment; filename=\"qrcode.png\""}
	imagePart["Content-Transfer-Encoding"] = []string{"base64"}
	imagePartBody := strings.NewReader(qrCodeBase64)

	// Create multipart message.
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for key, values := range htmlPart {
		for _, value := range values {
			writer.WriteField(key, value)
		}
	}
	htmlPartWriter, _ := writer.CreatePart(htmlPart)
	io.Copy(htmlPartWriter, htmlPartBody)
	imagePartWriter, _ := writer.CreatePart(imagePart)
	io.Copy(imagePartWriter, imagePartBody)
	writer.Close()

	// Define email message.
	from := m.From

	to := []string{register.Login}
	subject := "QR Code Test"
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/related; boundary=%s\r\n\r\n%s", from, strings.Join(to, ","), subject, writer.Boundary(), body))

	// Send email.
	err = smtp.SendMail(m.Host+":"+m.Port, smtp.PlainAuth("", m.Username, m.Password, m.Host), from, to, msg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Email sent!")
}
*/
