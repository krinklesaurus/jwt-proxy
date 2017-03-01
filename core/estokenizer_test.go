package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/SermoDigital/jose/crypto"
	"github.com/krinklesaurus/jwt_proxy"
	"github.com/krinklesaurus/jwt_proxy/mock"
	"github.com/stretchr/testify/assert"
)

func TestES(t *testing.T) {
	signingMethod := app.SigningMethods["ES256"]

	privKeyFile := "../mock/files/ec256-private.pem"

	privateKey, err := ioutil.ReadFile(privKeyFile)
	assert.Nil(t, err, fmt.Sprintf("could not read %s", privKeyFile))

	privateKeyString, err := crypto.ParseECPrivateKeyFromPEM(privateKey)
	assert.Nil(t, err, fmt.Sprintf("could not parse %s", privKeyFile))

	tokenizer := NewESTokenizer(signingMethod, privateKeyString)

	token, err := tokenizer.Serialize(mock.MockClaims)
	assert.Nil(t, err, "ec serializing failed")

	pubKeyFile := "../mock/files/ec256-public.pem"
	pubKeyPem, _ := ioutil.ReadFile(pubKeyFile)
	pubKey, _ := crypto.ParseECPublicKeyFromPEM(pubKeyPem)

	jsonString, _ := json.Marshal(mock.MockClaims)
	signingMethod.Verify(jsonString, token, pubKey)
}
