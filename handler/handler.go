package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alecthomas/template"
	"github.com/gorilla/mux"
	app "github.com/krinklesaurus/jwt_proxy"
	"github.com/krinklesaurus/jwt_proxy/log"
)

func New(core app.CoreAuth, nonceStore app.NonceStore) (*Handler, error) {
	return &Handler{core: core, nonceStore: nonceStore}, nil
}

type Handler struct {
	core       app.CoreAuth
	nonceStore app.NonceStore
}

func (handler *Handler) jwtHandler(w http.ResponseWriter, r *http.Request, token *app.Token) {
	claims := handler.core.Claims(token)
	tokenByte, err := handler.core.JwtToken(claims)
	if err != nil {
		log.Errorf("error %s", err.Error())
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	jwtAsString := string(tokenByte)

	url := handler.core.RedirectURI()
	urlWithToken := fmt.Sprintf(url+"?token=%s", jwtAsString)
	http.Redirect(w, r, urlWithToken, 302)
}

func (handler *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", 302)
}

func (handler *Handler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	if provider == "" {
		log.Errorf("missing provider param")
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	queryParams := r.URL.Query()
	code := queryParams.Get("code")
	state := queryParams.Get("state")

	nonce, err := handler.nonceStore.GetAndRemove(r)
	if err != nil {
		log.Errorf("Could not retrieve nonce from store %s", err.Error())
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	if code == "" || nonce != state {
		log.Errorf("missing code %s or states don't match: session:%s vs. param:%s", code, nonce, state)
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	token, err := handler.core.Token(provider, code)
	if err != nil {
		log.Errorf("error retrieving token %s", err.Error())
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	handler.jwtHandler(w, r, token)
}

func (handler *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginTemplate, err := template.ParseFiles("www/login.html")
	if err != nil {
		log.Errorf("error parsing www/login.html %s", err.Error())
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	supportedProviders := handler.core.Providers()

	csrf, err := handler.nonceStore.CreateNonce(w, r)
	if err != nil {
		log.Errorf("error creating csrf %s", err.Error())
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	templateData := struct {
		LocalAuthURL string
		Providers    []string
		CSRF         string
	}{
		"/auth",
		supportedProviders,
		csrf,
	}

	loginTemplate.Execute(w, templateData)
}

func (handler *Handler) ProviderLoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	state, err := handler.nonceStore.CreateNonce(w, r)
	if err != nil {
		log.Errorf("error creating nonce %s", err.Error())
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	authCodeURL := handler.core.AuthURL(provider, state)
	http.Redirect(w, r, authCodeURL, 302)
}

func (handler *Handler) PublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	publicKey, err := handler.core.PublicKey()
	if err != nil {
		log.Errorf("error reading public key %s", err.Error())
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(publicKey)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
