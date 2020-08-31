package core

import (
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/SermoDigital/jose/jws"
	app "github.com/krinklesaurus/jwt-proxy"
	"github.com/krinklesaurus/jwt-proxy/log"
	"golang.org/x/oauth2"
)

func New(config *app.Config, tokenizer app.Tokenizer, userService app.UserService) *Core {
	tokenStore := map[string]*app.TokenInfo{}
	return &Core{Config: config, userService: userService, tokenStore: tokenStore, Tokenizer: tokenizer}
}

type Core struct {
	Config      *app.Config
	userService app.UserService
	tokenStore  map[string]*app.TokenInfo
	Tokenizer   app.Tokenizer
}

func (c *Core) PublicKey() (*app.PublicKey, error) {
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

	return &app.PublicKey{Keys: keys}, nil
}

func (c *Core) TokenInfo(providerID string, code string) (*app.TokenInfo, error) {
	provider := c.Config.Providers[providerID]
	log.Debugf("getting access token from %s", provider.Name())
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

	token := &app.TokenInfo{Token: *providerToken, User: user, Provider: provider}
	return token, nil
}

func (c *Core) Claims(token *app.TokenInfo) (jws.Claims, error) {
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

func (c *Core) AuthURL(providerID string, state string) string {
	provider := c.Config.Providers[providerID]
	url := provider.AuthCodeURL(state)
	return url
}
