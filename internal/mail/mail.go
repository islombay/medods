package mail

import (
	"fmt"
	"medods/config"
	"medods/internal/model"
	"medods/pkg/logs"
	"net/smtp"
	"os"
)

type MailService struct {
	log logs.LoggerInterface
	cfg config.ConfigMail
}

type MailServiceI interface {
	Send(to, message string) error
	WarnIPAddressChange(to string, info model.DeviceInfo) error
}

func NewMailService(log logs.LoggerInterface, cfg config.ConfigMail) MailServiceI {
	return &MailService{
		log: log,
		cfg: cfg,
	}
}

func (mail *MailService) Send(to, message string) error {
	mail.log.Debug("sending email", logs.String("to", to))
	mail.log.Debug("credentials for mail",
		logs.String("pwd", os.Getenv("MAIL_PWD")),
		logs.String("host", mail.cfg.Host),
		logs.String("port", mail.cfg.Port),
	)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", mail.cfg.Host, mail.cfg.Port),
		smtp.PlainAuth("", mail.cfg.Sender, os.Getenv("MAIL_PWD"), mail.cfg.Host),
		mail.cfg.Sender,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		mail.log.Error("could not send email", logs.Error(err),
			logs.String("to", to))
		return err
	}

	mail.log.Debug("email sent", logs.String("to", to))

	return nil
}

func (mail *MailService) WarnIPAddressChange(to string, info model.DeviceInfo) error {
	msg := fmt.Sprintf("Новый вход с нового IP адреса.\nIP адрес: %s, User-Agent: %s", info.IP, info.UserAgent)

	return mail.Send(to, msg)
}
