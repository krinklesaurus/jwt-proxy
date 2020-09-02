package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SermoDigital/jose/jws"

	"github.com/alecthomas/template"
	"github.com/gorilla/mux"
	"github.com/krinklesaurus/jwt-proxy/config"
	"github.com/krinklesaurus/jwt-proxy/core"
	"github.com/krinklesaurus/jwt-proxy/log"
)

type Handler struct {
	config     *config.Config
	core       core.CoreAuth
	nonceStore NonceStore
}

// PublicKey is a struct for a list of keys
type PublicKey struct {
	Keys []string
}

func New(config *config.Config, core core.CoreAuth, nonceStore NonceStore) (*Handler, error) {
	return &Handler{config: config, core: core, nonceStore: nonceStore}, nil
}

func (handler *Handler) jwtHandler(w http.ResponseWriter, r *http.Request, token *core.TokenInfo) {
	claims, err := handler.core.Claims(token)
	if err != nil {
		log.Errorf("error %s", err.Error())
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

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
	http.Redirect(w, r, "/jwt-proxy/login", 302)
}

func (handler *Handler) RobotsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("%s/%s", handler.config.WWWRootDir, "robots.txt"))
}

func (handler *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (handler *Handler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	providerName := vars["provider"]
	if providerName == "" {
		log.Errorf("missing provider param")
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	queryParams := r.URL.Query()
	code := queryParams.Get("code")
	state := queryParams.Get("state")

	log.Debugf("received code %s and state %s", code, state)

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

	token, err := handler.core.GenTokenInfo(providerName, code)
	if err != nil {
		log.Errorf("error retrieving token %s", err.Error())
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	handler.jwtHandler(w, r, token)
}

func (handler *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginTemplate, err := template.ParseFiles(fmt.Sprintf("%s/%s", handler.config.WWWRootDir, "login.html"))
	if err != nil {
		log.Errorf("error parsing %s, error is %v", fmt.Sprintf("%s/%s", handler.config.WWWRootDir, "login.html"), err.Error())
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
	log.Debugf("redirecting to %s", authCodeURL)
	http.Redirect(w, r, authCodeURL, 302)
}

func (handler *Handler) PublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	publicKeys, err := handler.core.PublicKeys()
	if err != nil {
		log.Errorf("error reading public key %s", err.Error())
		http.Error(w, "Sorry, some unknown error occurred", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(struct {
		Keys []string `json:"keys"`
	}{
		Keys: publicKeys,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (handler *Handler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	jwt, err := jws.ParseJWTFromRequest(r)
	if err != nil {
		log.Errorf("no jwt found: %v", err)
		http.Error(w, "no jwt found", http.StatusUnauthorized)
		return
	}
	err = jwt.Validate(handler.config.PublicRSAKey, core.SigningMethods[handler.config.SigningMethod])
	if err != nil {
		log.Errorf("no valid jwt: %v", err)
		http.Error(w, "no valid jwt", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
