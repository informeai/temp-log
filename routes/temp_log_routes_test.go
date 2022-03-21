package routes

import (
	"testing"
)

//go test -run ^TestHealth
func TestHealth(t *testing.T) {
	r := NewRouter()
	if err := r.health(); err != nil {
		t.Errorf("TestHealth: got: %s -> expect: == nil\n", err.Error())
	}
}

//go test -run ^TestGetLogs
func TestGetLogs(t *testing.T) {
	r := NewRouter()
	if err := r.getLogs(); err != nil {
		t.Errorf("TestGetLogs: got: %s -> expect: == nil\n", err.Error())
	}
}

//go test -run ^TestPostLogs
func TestPostLogs(t *testing.T) {
	r := NewRouter()
	if err := r.postLogs(); err != nil {
		t.Errorf("TestPostLogs: got: %s -> expect: == nil\n", err.Error())
	}
}

//go test -run ^TestAuth
func TestAuth(t *testing.T) {
	r := NewRouter()
	if err := r.auth(); err != nil {
		t.Errorf("TestAuth: got: %s -> expect: == nil\n", err.Error())
	}
}
