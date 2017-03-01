package core

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/krinklesaurus/jwt_proxy"
	"github.com/krinklesaurus/jwt_proxy/mock"
	"github.com/stretchr/testify/assert"
)

func TestRSA(t *testing.T) {
	signingMethod := app.SigningMethods["RS256"]

	privKeyFile := "../mock/files/sample_key.priv"

	privateKey, err := ioutil.ReadFile(privKeyFile)
	assert.Nil(t, err, fmt.Sprintf("could not read %s", privKeyFile))

	privateKeyString, err := crypto.ParseRSAPrivateKeyFromPEM(privateKey)
	assert.Nil(t, err, fmt.Sprintf("could not parse %s", privKeyFile))

	tokenizer := NewRSATokenizer(signingMethod, privateKeyString)

	token, err := tokenizer.Serialize(mock.MockClaims)
	assert.Nil(t, err, "rsa serializing failed")

	pubKeyFile := "../mock/files/sample_key.pub"
	pubKeyPem, _ := ioutil.ReadFile(pubKeyFile)
	pubKey, _ := crypto.ParseRSAPublicKeyFromPEM(pubKeyPem)

	jwt, err := jws.ParseJWT(token)
	assert.Nil(t, err, "rsa jwt parsing failed")

	err = jwt.Validate(pubKey, signingMethod)
	assert.Nil(t, err, "rsa verifying failed")
}
