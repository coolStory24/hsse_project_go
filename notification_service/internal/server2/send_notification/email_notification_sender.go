package send_notification

import (
	"fmt"
	"hotel_service/internal/models"
	"hotel_service/internal/server2/content_build"
	"net/smtp"
)

type EmailSender struct {
	SMTPServer string
	Port       int
	Username   string
	Password   string
}

func (e *EmailSender) Send(booking models.BookingRequest, contentBuilder content_build.IContentBuilder, clientData models.ClientData) error {
	content := contentBuilder.BuildContent(booking)

	auth := smtp.PlainAuth("", e.Username, e.Password, e.SMTPServer)

	from := e.Username
	to := []string{clientData.Email}
	subject := "Booking Confirmation"
	body := content

	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	err := smtp.SendMail(fmt.Sprintf("%s:%d", e.SMTPServer, e.Port), auth, from, to, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
