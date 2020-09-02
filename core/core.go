package core

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/krinklesaurus/jwt-proxy/config"
	"github.com/krinklesaurus/jwt-proxy/log"
	"github.com/krinklesaurus/jwt-proxy/provider"
	"github.com/krinklesaurus/jwt-proxy/user"
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

// CoreAuth is the central interface of jwt-proxy. It provides all function necessary
// for handling the redirect to the provider, process the login, enrich the provider's
// token with some custom parameters and return the JWT token to the callback URI.
type CoreAuth interface {
	PublicKeys() ([]string, error)
	GenTokenInfo(provider string, code string) (*TokenInfo, error)
	Claims(token *TokenInfo) (jws.Claims, error)
	JwtToken(jws.Claims) ([]byte, error)
	RedirectURI() string
	AuthURL(provider string, state string) (string, error)
	Providers() []string
}

// TokenInfo wraps oauth.Token and adds two additional fields:
// Provider is the OAuth provider, e.g. github or facebook
// UserInfo is the map of user info claims from the provider
type TokenInfo struct {
	oauth2.Token
	Provider provider.Provider
	User     string
}

func New(config *config.Config, tokenizer Tokenizer, userService user.UserService) *Core {
	tokenStore := map[string]*TokenInfo{}
	return &Core{Config: config, userService: userService, tokenStore: tokenStore, Tokenizer: tokenizer}
}

type Core struct {
	Config      *config.Config
	userService user.UserService
	tokenStore  map[string]*TokenInfo
	Tokenizer   Tokenizer
}

func (c *Core) PublicKeys() ([]string, error) {
	publicKey := c.Config.PrivateRSAKey.PublicKey
	publicKeyDer, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, err
	}

	publicKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	publicKeyPem := string(pem.EncodeToMemory(&publicKeyBlock))

	keys := []string{
		publicKeyPem,
	}

	return keys, nil
}

func (c *Core) GenTokenInfo(providerID string, code string) (*TokenInfo, error) {
	provider := c.Config.Providers[providerID]
	log.Debugf("getting access token from %s with code %s", provider.Name(), code)
	providerToken, err := provider.Exchange(oauth2.NoContext, code)

	if err != nil {
		return nil, err
	}

	log.Debugf("received provider token %+v", providerToken)

	userID, err := provider.User()
	if err != nil {
		return nil, err
	}
	user, err := c.userService.UniqueUser(providerID, userID)
	if err != nil {
		return nil, err
	}

	token := &TokenInfo{Token: *providerToken, User: user, Provider: provider}
	return token, nil
}

func (c *Core) Claims(token *TokenInfo) (jws.Claims, error) {
	log.Debugf("received token %s from provider %s", token.AccessToken, token.Provider.Name())
	// see https://openid.net/specs/openid-connect-core-1_0.html#IDToken

	claims := jws.Claims{}
	claims.SetIssuer(c.Config.RootURI)
	claims.SetSubject(c.Config.Subject)
	claims.SetAudience(token.Provider.ClientID())

	now := time.Now()
	expiry := token.Expiry
	if aft := expiry.After(now); !aft {
		expiry = now.Add(time.Duration(c.Config.ExpirySeconds) * time.Second)
	}
	claims.SetExpiration(expiry)
	claims.SetIssuedAt(time.Now())

	claims.Set("provider", token.Provider.Name())
	claims.Set("user", token.User)
	claims.Set("access_token", token.AccessToken)
	claims.Set("token_type", token.TokenType)
	claims.Set("refresh_token", token.RefreshToken)

	return claims, nil
}

func (c *Core) JwtToken(claims jws.Claims) ([]byte, error) {
	b, err := c.Tokenizer.Serialize(claims)

	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Core) RedirectURI() string {
	return c.Config.RedirectURI
}

func (c *Core) LocalEnabled() bool {
	return false
}

func (c *Core) Providers() []string {
	keys := []string{}
	for key := range c.Config.Providers {
		keys = append(keys, key)
	}
	return keys
}

func (c *Core) AuthURL(providerID string, state string) (string, error) {
	provider := c.Config.Providers[providerID]
	if provider == nil {
		return "", fmt.Errorf("provider %s not found", providerID)
	}
	url := provider.AuthCodeURL(state)
	return url, nil
}
