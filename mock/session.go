package mock

import (
	"net/http"

	"github.com/krinklesaurus/jwt-proxy"
)

func NewSessionStore() app.NonceStore {
	return &sessionStore{}
}

type sessionStore struct {
	nonce string
}

func (store *sessionStore) CreateNonce(w http.ResponseWriter, r *http.Request) (string, error) {
	return "csrf", nil
}

func (store *sessionStore) GetAndRemove(r *http.Request) (string, error) {
	return "csrf", nil
}
