package mailer

import (
	"github.com/wneessen/go-mail"
)

func SendEmail(from, to string, msg mail.Msg) error {
	if err := msg.From(from); err != nil {
		return err
	}
	if err := msg.To(to); err != nil {
		return err
	}
	return nil
}
