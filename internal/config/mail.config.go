package config

import (
	"os"
)

type MailConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	SenderEmail string
}

func LoadMailConfig() MailConfig {
	return MailConfig{
		Host:        "smtp.gmail.com",
		Port:        587,
		Username:    os.Getenv("SMTP_USER"),
		Password:    os.Getenv("SMTP_PASS"),
		SenderEmail: os.Getenv("SMTP_SENDER_EMAIL"),
	}
}
