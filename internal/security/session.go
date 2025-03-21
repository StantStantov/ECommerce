package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

func NewSessionCookie() (*http.Cookie, error) {
	token, err := generateToken(64)
	if err != nil {
		return nil, fmt.Errorf("security.NewSessionCookie: [%w]", err)
	}
	return &http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: true,
	}, nil
}

func NewCsrfCookie() (*http.Cookie, error) {
	token, err := generateToken(64)
	if err != nil {
		return nil, fmt.Errorf("security.NewCsrfCookie: [%w]", err)
	}
	return &http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: false,
	}, nil
}

func generateToken(length int) (string, error) {
	token := make([]byte, length)
	if _, err := rand.Read(token); err != nil {
		return "", fmt.Errorf("security.generateToken: [%w]", err)
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
