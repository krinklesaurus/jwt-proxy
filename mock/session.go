package mock

import "net/http"

func NewMockSessionStore() *MockSessionStore {
	return &MockSessionStore{}
}

type MockSessionStore struct {
	nonce string
}

func (store *MockSessionStore) CreateNonce(w http.ResponseWriter, r *http.Request) (string, error) {
	return "csrf", nil
}

func (store *MockSessionStore) GetAndRemove(r *http.Request) (string, error) {
	return "csrf", nil
}
