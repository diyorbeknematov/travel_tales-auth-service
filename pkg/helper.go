package pkg

import "net/smtp"

func SendEmail(to string, subject string, body string) error {
	from := "nematovdiyorbek10@gmail.com"
	password := "software_engineer"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

// Create a reset link
func CreateResetLink(baseURL string, token string) string {
	return baseURL + "new-password?token=" + token
}
