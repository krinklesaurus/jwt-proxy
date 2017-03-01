package handler

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/krinklesaurus/jwt_proxy/log"
	"github.com/krinklesaurus/jwt_proxy/util"
)

const sessionName string = "nonce-session"
const sessionNonce string = "nonce"

func NewHttpSessionStore() (*HttpSessionStore, error) {
	sessionStore := sessions.NewCookieStore([]byte(util.RandomString(32)))
	return &HttpSessionStore{sessionStore: sessionStore}, nil
}

type HttpSessionStore struct {
	sessionStore *sessions.CookieStore
}

func (store *HttpSessionStore) CreateNonce(w http.ResponseWriter, r *http.Request) (string, error) {
	nonce := util.RandomString(32)
	log.Debugf("get sessions store with name %s", sessionName)
	session, err := store.sessionStore.Get(r, sessionName)
	if err != nil {
		log.Errorf("error getting session: %v", err)
	}
	log.Debugf("set session key %s to %s", sessionNonce, nonce)
	session.Values[sessionNonce] = nonce
	err = session.Save(r, w)
	if err != nil {
		log.Errorf("error saving session: %v", err)
	}
	return nonce, nil
}

func (store *HttpSessionStore) GetAndRemove(r *http.Request) (string, error) {
	session, err := store.sessionStore.Get(r, sessionName)
	if err != nil {
		return "", err
	}
	value, ok := session.Values[sessionNonce].(string)
	session.Values[sessionNonce] = nil
	if !ok {
		return "", errors.New("value from session is not a string")
	}
	return value, nil
}
