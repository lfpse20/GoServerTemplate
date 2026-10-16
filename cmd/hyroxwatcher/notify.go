package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

// sendAlert sends an email (or SMS via carrier gateway) alert to all configured recipients.
func sendAlert(cfg Config, subject, body string) error {
	if cfg.GmailUser == "" || cfg.GmailAppPassword == "" {
		return fmt.Errorf("GMAIL_USER or GMAIL_APP_PASSWORD not set")
	}
	if len(cfg.NotifyEmails) == 0 {
		return fmt.Errorf("NOTIFY_EMAILS not set")
	}

	auth := smtp.PlainAuth("", cfg.GmailUser, cfg.GmailAppPassword, "smtp.gmail.com")

	msg := []byte(
		"From: " + cfg.GmailUser + "\r\n" +
			"To: " + strings.Join(cfg.NotifyEmails, ", ") + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	return smtp.SendMail("smtp.gmail.com:587", auth, cfg.GmailUser, cfg.NotifyEmails, msg)
}
