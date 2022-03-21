package services

import "testing"

//go test -run ^TestVerifyEmail
func TestVerifyEmail(t *testing.T) {
	email := "contato.informeai@gmail.com"
	if err := VerifyEmail(email); err != nil {
		t.Errorf("TestVerifyEmail: got -> %s , expect -> nil", err.Error())
	}
}
