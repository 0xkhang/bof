package utils

import "net/mail"

func IsEmailValid(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	return nil
}
