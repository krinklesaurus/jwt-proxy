package core

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/SermoDigital/jose/crypto"
	"github.com/krinklesaurus/jwt-proxy/config"
	"github.com/krinklesaurus/jwt-proxy/user"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

type mockProvider struct {
	userId string
}

func (m mockProvider) AuthCodeURL(state string) string {
	return ""
}

func (m mockProvider) User() (string, error) {
	return m.userId, nil
}

func (m mockProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return &oauth2.Token{}, nil
}

func (m mockProvider) String() string {
	return m.Name()
}

func (m mockProvider) Name() string {
	return "mock_provider"
}

func (m mockProvider) ClientID() string {
	return "mock_clientid"
}

type mockUserservice struct {
}

func TestPublicKey(t *testing.T) {
	conf, _ := config.Initialize("../config-test.yml")
	core := New(conf, nil, nil)
	publicKeys, _ := core.PublicKeys()

	for _, value := range publicKeys {
		block, _ := pem.Decode([]byte(value))
		assert.NotNil(t, block, "block should be not be null")

		_, err := x509.ParsePKIXPublicKey([]byte(block.Bytes))
		assert.Nil(t, err, "could not parse public key")
	}
}

func TestJwtToken(t *testing.T) {
	conf, _ := config.Initialize("../config-test.yml")
	userID := uuid.NewV4().String()
	conf.Providers["mock_provider"] = mockProvider{userId: userID}

	core := New(conf, NewRSATokenizer(crypto.SigningMethodRS256, conf.PrivateRSAKey), user.PlainUserService{})
	token, err := core.GenTokenInfo("mock_provider", "code")
	assert.Nil(t, err, "err should be nothing")

	claims, _ := core.Claims(token)
	assert.Nil(t, err, "err should be nothing")

	data, err := core.JwtToken(claims)
	assert.Nil(t, err, "err should be nothing")

	assert.NotEmpty(t, data)
}
