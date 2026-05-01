package service

import (
	"alumnihub/internal/config"
	"fmt"
	"net/smtp"
)

type EmailService struct {
	Config config.AppConfig
}

func NewEmailService(cfg config.AppConfig) *EmailService {
	return &EmailService{
		Config: cfg,
	}
}

type EmailData struct {
	To      string
	Subject string
	Body    string
}

// SendEmail sends an email using SMTP
func (e *EmailService) SendEmail(data EmailData) error {
	if e.Config.SMTPHost == "" || e.Config.SMTPUsername == "" || e.Config.SMTPPassword == "" || e.Config.SMTPFrom == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	// Set up authentication
	auth := smtp.PlainAuth("", e.Config.SMTPUsername, e.Config.SMTPPassword, e.Config.SMTPHost)

	// Construct the email with required headers
	to := []string{data.To}
	from := e.Config.SMTPFrom
	if from == "" {
		from = e.Config.SMTPUsername
	}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		from,
		data.To,
		data.Subject,
		data.Body,
	))

	// Send the email
	addr := fmt.Sprintf("%s:%d", e.Config.SMTPHost, e.Config.SMTPPort)
	err := smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email
func (e *EmailService) SendPasswordResetEmail(to, resetToken string) error {
	resetURL := fmt.Sprintf("http://%s/reset_password?token=%s", e.Config.Domain, resetToken)

	subject := "Password Reset Request - AlumniHub"

	body := fmt.Sprintf(`Hello,

You have requested to reset your password for your AlumniHub account.

Please click the following link to reset your password:
%s

This link will expire in 30 minutes for security reasons.

If you did not request this password reset, please ignore this email.

Best regards,
AlumniHub Team

---
This is an automated message. Please do not reply to this email.`, resetURL)

	return e.SendEmail(EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
	})
}

// SendPasswordResetSuccessEmail sends a confirmation email after password reset
func (e *EmailService) SendPasswordResetSuccessEmail(to string) error {
	subject := "Password Reset Successful - AlumniHub"

	body := `Hello,

Your password has been successfully reset for your AlumniHub account.

If you did not perform this action, please contact our support team immediately.

Best regards,
AlumniHub Team

---
This is an automated message. Please do not reply to this email.`

	return e.SendEmail(EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
	})
}
