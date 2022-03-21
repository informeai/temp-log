package services

import (
	"errors"
	"strings"
)

func VerifyEmail(email string) error {
	if len(email) < 10 {
		return errors.New("not valid email")
	}
	existSymbol := strings.Contains(email, "@")
	if existSymbol == false {
		return errors.New("not valid email")
	}
	indexSymbol := strings.Index(email, "@")
	indexDot := strings.LastIndex(email, ".")
	if indexDot == -1 || indexDot < indexSymbol {
		return errors.New("not valid email")
	}
	return nil
}
