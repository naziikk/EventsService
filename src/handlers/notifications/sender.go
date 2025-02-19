package notifications

import (
	"RedisService/internal/config"
	"log"
	"net/smtp"
)

func SendEmail(recipient string, subject string, body string) error {
	cfg := config.MustLoadConfig()
	message := []byte("Subject: " + subject + "\r\n" + body)
	auth := smtp.PlainAuth("", cfg.YandexSMTP.SenderEmail, cfg.YandexSMTP.AuthApiKey, cfg.YandexSMTP.SMTPServer)
	err := smtp.SendMail(cfg.YandexSMTP.SMTPServer+":"+cfg.YandexSMTP.SMTPPort, auth, cfg.YandexSMTP.SenderEmail, []string{recipient}, message)
	if err != nil {
		log.Printf("Ошибка при отправке письма: %v", err)
		return err
	}
	log.Printf("Письмо отправлено на адрес %s", recipient)
	return nil
}
