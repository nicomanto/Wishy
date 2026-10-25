package controllers

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestManagerPasswordVerification(t *testing.T) {
	password := "generated-password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	if !verifyManagerPassword(string(hash), password) {
		t.Fatal("verifyManagerPassword() rejected the correct password")
	}
	if verifyManagerPassword(string(hash), password+"wrong") {
		t.Fatal("verifyManagerPassword() accepted an incorrect password")
	}
	if verifyManagerPassword(password, password) {
		t.Fatal("verifyManagerPassword() accepted a plaintext password")
	}
}
