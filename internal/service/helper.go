package service

import (
	"net/mail"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

func validEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func hasUpperAndLower(s string) bool {
	return s != strings.ToLower(s) && s != strings.ToUpper(s)
}

func strongPassword(password string) bool {

	const specialSymbols = "!@#$%^&*()_+-={}[]:;\"'<>,.?/\\|`~"
	const digits = "0123456789"

	if utf8.RuneCountInString(password) < 8 {
		return false
	}

	if !strings.ContainsAny(password, specialSymbols) {
		return false
	}

	if !strings.ContainsAny(password, digits) {
		return false
	}

	if !hasUpperAndLower(password) {
		return false
	}
	return true
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
