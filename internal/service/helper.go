package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/mail"
	"strings"
	"unicode"
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

func generateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func validName(name string) bool {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return false
	}

	for _, r := range trimmed {
		if unicode.IsControl(r) {
			return false
		}
	}

	return true
}
