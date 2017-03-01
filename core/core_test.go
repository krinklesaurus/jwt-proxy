package core

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/krinklesaurus/jwt_proxy"
	"github.com/krinklesaurus/jwt_proxy/mock"
	"github.com/stretchr/testify/assert"
)

var mockConfig = mock.NewMockConfig()
var mockUserService = mock.NewUserservice()
var mockTokenizer = mock.NewTokenizer()
var core app.CoreAuth = New(mockConfig, mockTokenizer, mockUserService)

func TestPublicKey(t *testing.T) {
	publicKey, _ := core.PublicKey()

	keys := publicKey.Keys
	for _, value := range keys {
		block, _ := pem.Decode([]byte(value))
		assert.NotNil(t, block, "block should be not be null")

		_, err := x509.ParsePKIXPublicKey([]byte(block.Bytes))
		assert.Nil(t, err, "could not parse public key")
	}
}

func TestToken(t *testing.T) {
	token, err := core.Token("MOCK_PROVIDER", "code")
	assert.Nil(t, err, "err should be nothing")

	assert.Equal(t, "MOCK_PROVIDER", token.ProviderID, "provider ids do not match")
	assert.Equal(t, "MOCK_USER_ID", token.UserID, "user ids do not match")
}

func TestJwtToken(t *testing.T) {
	token, err := core.Token("MOCK_PROVIDER", "code")
	assert.Nil(t, err, "err should be nothing")

	claims := core.Claims(token)
	assert.Nil(t, err, "err should be nothing")

	data, err := core.JwtToken(claims)
	assert.Nil(t, err, "err should be nothing")

	assert.Equal(t, "MOCK_JWT_TOKEN", string(data), "jwt token do not match")
}
