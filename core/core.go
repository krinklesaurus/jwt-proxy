package core

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/krinklesaurus/jwt_proxy"
	"github.com/krinklesaurus/jwt_proxy/log"
	"github.com/krinklesaurus/jwt_proxy/util"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
)

func New(config *app.Config, tokenizer app.Tokenizer, userService app.UserService) *Core {
	tokenStore := map[string]*app.Token{}
	return &Core{Config: config, userService: userService, tokenStore: tokenStore, Tokenizer: tokenizer}
}

type Core struct {
	Config      *app.Config
	userService app.UserService
	tokenStore  map[string]*app.Token
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

func (c *Core) Token(providerID string, code string) (*app.Token, error) {
	provider := c.Config.Providers[providerID]
	providerToken, err := provider.Exchange(oauth2.NoContext, code)

	if err != nil {
		return nil, err
	}

	log.Debugf("received provider token AccessToken: %s RefreshToken: %s TokenType: %s Expiry: %v",
		providerToken.AccessToken, providerToken.RefreshToken, providerToken.TokenType, providerToken.Expiry)

	providerUserID, err := provider.UniqueUserID()
	if err != nil {
		return nil, err
	}

	user, err := c.userService.LoadUser(providerID, providerUserID)
	if err != nil {
		return nil, err
	}

	token := &app.Token{Token: *providerToken, ProviderID: providerID, UserID: user.ID}
	return token, nil
}

func (c *Core) Claims(token *app.Token) jws.Claims {
	log.Debugf("received token %s from provider %s", token.AccessToken, token.ProviderID)

	claims := jws.Claims{}
	claims.SetAudience(c.Config.Audience)
	claims.SetExpiration(token.Expiry)
	claims.SetIssuedAt(time.Now())
	claims.SetIssuer(c.Config.Issuer)
	claims.SetJWTID(uuid.NewV4().String())
	claims.SetNotBefore(time.Now())
	claims.SetSubject(c.Config.Subject)

	claims.Set("provider", token.ProviderID)
	claims.Set("user", token.UserID)
	claims.Set("access_token", token.AccessToken)
	claims.Set("token_type", token.TokenType)
	claims.Set("refresh_token", token.RefreshToken)
	claims.Set("expiry", token.Expiry.Unix())

	return claims
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

func (c *Core) Auth(username string, password string) (*app.Token, error) {
	err := c.userService.LoginUser(username, password)
	if err != nil {
		log.Errorf("Failed to login user %s cause of %s", username, err)
		return nil, errors.New("Wrong credentials")
	}

	tokenValue := util.RandomString(32)
	token := &app.Token{Token: oauth2.Token{
		AccessToken:  tokenValue,
		TokenType:    "bearer",
		RefreshToken: "",
		Expiry:       time.Now(),
	}, ProviderID: "local"}
	c.tokenStore[username] = token
	log.Infof("added to token store %s => %s", username, tokenValue)
	return token, nil
}
