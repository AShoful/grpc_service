package handler

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "testsecret")
	code := m.Run()
	os.Exit(code)
}
