package send_notification

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"notification_service/internal/models"
	"notification_service/internal/service/content_build"
)

type EmailSender struct {
	contentBuilder content_build.IContentBuilder
	SMTPServer     string
	Port           int
	Username       string
	Password       string
}

func NewEmailSender(contentBuilder content_build.IContentBuilder, SMTPServer string, port int, username string, password string) *EmailSender {
	return &EmailSender{
		contentBuilder: contentBuilder,
		SMTPServer:     SMTPServer,
		Port:           port,
		Username:       username,
		Password:       password,
	}
}

func (e *EmailSender) Send(notification models.NotificationData) error {
	content := e.contentBuilder.BuildContent(notification)

	auth := smtp.PlainAuth("", e.Username, e.Password, e.SMTPServer)

	from := e.Username
	to := []string{notification.UserContactData.Email}
	subject := "Booking Confirmation"
	body := content

	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	slog.Info("Sending email...")

	_ = smtp.SendMail(fmt.Sprintf("%s:%d", e.SMTPServer, e.Port), auth, from, to, []byte(message))

	//if err != nil {
	//	slog.Error("Failed to send email")
	//}

	// Explanation: SendMail will always throw error, because we don't have a registered email address and SMTP
	// from which we can send emails.
	// If we had it, we could put the data inti .env file and the commented out code above would work
	slog.Info("\"AS IF\" Email successfully sent")

	return nil
}
