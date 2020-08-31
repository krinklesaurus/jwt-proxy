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

// TokenInfo wraps oauth.Token and adds two additional fields:
// Provider is the OAuth provider, e.g. github or facebook
// UserInfo is the map of user info claims from the provider
type TokenInfo struct {
	oauth2.Token
	Provider Provider
	User     string
}

// PublicKey is a struct for a list of keys
type PublicKey struct {
	Keys []string
}

// CoreAuth is the central interface of jwt-proxy. It provides all function necessary
// for handling the redirect to the provider, process the login, enrich the provider's
// token with some custom parameters and return the JWT token to the callback URI.
type CoreAuth interface {
	PublicKey() (*PublicKey, error)
	TokenInfo(provider string, code string) (*TokenInfo, error)
	Claims(token *TokenInfo) (jws.Claims, error)
	JwtToken(jws.Claims) ([]byte, error)
	RedirectURI() string
	AuthURL(provider string, state string) string
	Providers() []string
}

// NonceStore simply stores a nonce for CSRF attack prevention
type NonceStore interface {
	CreateNonce(w http.ResponseWriter, r *http.Request) (string, error)
	GetAndRemove(r *http.Request) (string, error)
}

// UserService provides a function for creating a user from the given provider and providerUserID.
// The created user contains the global unique user ID that is used within your environment.
// user/hashuserservice is the most basic way to create a unique user by simply
// hashing both the provider and providerUserID. It would also be possible to load
// the user id from a DB using the provider and providerUserID as a key or just
// concatenate both strings.
type UserService interface {
	UniqueUser(provider string, providerUserID string) (string, error)
}

// Provider is the interface every OAuth provider has to fulfill for being
// used within jwt-proxy
type Provider interface {
	AuthCodeURL(state string) string
	User() (string, error)
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	String() string
	Name() string
	ClientID() string
}

// Tokenizer creates a byte array from an input map.
type Tokenizer interface {
	Serialize(claims map[string]interface{}) ([]byte, error)
}

// Config is the struct for the config defined in a config.yml
type Config struct {
	RootURI       string
	RedirectURI   string
	WWWRootDir    string
	Providers     map[string]Provider
	SigningMethod crypto.SigningMethod
	PrivateRSAKey *rsa.PrivateKey
	PublicRSAKey  interface{}
	Audience      string
	Issuer        string
	Subject       string
	Password      string
	ExpirySeconds int
}

// String is a helping toString function for the config for debugging
func (c *Config) String() string {
	providersString := ""
	for _, p := range c.Providers {
		providersString = providersString + fmt.Sprintf(" %s with clientId %s ,", p.Name(), p.ClientID())
	}
	return fmt.Sprintf("rootURI: %s, redirectURI: %s, Audience: %s, Issuer: %s, Subject: %s, Providers: %s", c.RootURI, c.RedirectURI, c.Audience, c.Issuer, c.Subject, providersString)
}

// Log is an interface for logging
type Log interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
