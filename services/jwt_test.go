package services

import (
	"log"
	"testing"
)

//go test -run ^TestCreateJWT
func TestCreateJWT(t *testing.T) {
	email := "wellington@gmail.com"
	token, err := CreateJWT(email)
	if err != nil {
		t.Errorf("TestCreateJWT: got -> %s, expect -> nil", err.Error())
	}
	log.Println(token)
}

//go test -run ^TestVerifyJWT
func TestVerifyJWT(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IndlbGxpbmd0b25AZ21haWwuY29tIiwiZXhwIjoxNjQ3ODU5MzIzfQ.ZP4Tz2WXakjEN4uegXJKEQKkKxbLqeRe-uDW1XcW9rA"
	_, err := VerifyJWT(token)
	if err != nil {
		t.Errorf("TestVerifyJWT: got -> %s, expect -> nil", err.Error())
	}
}
