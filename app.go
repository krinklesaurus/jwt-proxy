package app

import (
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var SigningMethods = map[string]crypto.SigningMethod{
	"RS256": crypto.SigningMethodRS256,
	"RS384": crypto.SigningMethodRS384,
	"RS512": crypto.SigningMethodRS512,
	"HS256": crypto.SigningMethodHS256,
	"HS384": crypto.SigningMethodHS384,
	"HS512": crypto.SigningMethodHS512,
	"ES256": crypto.SigningMethodES256,
	"ES384": crypto.SigningMethodES384,
	"ES512": crypto.SigningMethodES512,
}

type Token struct {
	oauth2.Token
	ProviderID string `json:"provider"`
	UserID     string `json:"user"`
}

type PublicKey struct {
	Keys []string
}

type CoreAuth interface {
	PublicKey() (*PublicKey, error)
	Token(provider string, code string) (*Token, error)
	Claims(token *Token) jws.Claims
	JwtToken(jws.Claims) ([]byte, error)
	RedirectURI() string
	AuthURL(provider string, state string) string
	Providers() []string
	Auth(username string, password string) (*Token, error)
	LocalEnabled() bool
}

type NonceStore interface {
	CreateNonce(w http.ResponseWriter, r *http.Request) (string, error)
	GetAndRemove(r *http.Request) (string, error)
}

type User struct {
	ID string
}

type UserService interface {
	LoadUser(provider string, providerUserID string) (User, error)
	LoginUser(username string, plainPassword string) error
}

type Provider interface {
	AuthCodeURL(state string) string
	UniqueUserID() (string, error)
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
}

type Tokenizer interface {
	Serialize(claims map[string]interface{}) ([]byte, error)
}

type Config struct {
	RootURI       string
	RedirectURI   string
	Providers     map[string]Provider
	SigningMethod crypto.SigningMethod
	PrivateRSAKey *rsa.PrivateKey
	PublicRSAKey  interface{}
	Audience      string
	Issuer        string
	Subject       string
	LocalEnabled  bool
	Password      string
}

func (c *Config) String() string {
	return fmt.Sprintf("rootURI: %s, redirectURI: %s, Audience: %s, Issuer: %s, Subject: %s", c.RootURI, c.RedirectURI, c.Audience, c.Issuer, c.Subject)
}

type Log interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
