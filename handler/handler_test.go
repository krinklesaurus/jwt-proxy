package handler

import (
	"bytes"
	"testing"

	"net/http/httptest"
	"net/url"

	"github.com/krinklesaurus/jwt_proxy/mock"
)

func TestAuthHandler(t *testing.T) {
	mockCore := mock.NewCore()
	mockStore := mock.NewSessionStore()

	handler, _ := New(mockCore, mockStore)

	values := url.Values{}
	values.Add("username", "username")
	values.Add("password", "password")
	values.Add("csrf", "csrf")

	encodedValus := values.Encode()

	request := httptest.NewRequest("POST", "/auth", bytes.NewBufferString(encodedValus))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler.AuthHandler(w, request)
}
