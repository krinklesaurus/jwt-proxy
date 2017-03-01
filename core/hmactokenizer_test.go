package core

import (
	"io/ioutil"
	"testing"

	"github.com/SermoDigital/jose/jws"
	"github.com/krinklesaurus/jwt_proxy"
	"github.com/krinklesaurus/jwt_proxy/mock"
	"github.com/stretchr/testify/assert"
)

func TestHMAC(t *testing.T) {
	signingMethod := app.SigningMethods["HS256"]
	hmacKey := []byte("qwertyuiopasdfghjklzxcvbnm123456")

	_, err := ioutil.ReadFile("../mock/files/hmacTestKey")
	assert.Nil(t, err, "could not read mock/files/hmacTestKey")

	tokenizer := NewHMACTokenizer(signingMethod, hmacKey)

	token, err := tokenizer.Serialize(mock.MockClaims)
	assert.Nil(t, err, "hmac serializing failed")

	jwt, err := jws.ParseJWT(token)
	assert.Nil(t, err, "hmac jwt parsing failed")

	err = jwt.Validate(hmacKey, signingMethod)
	assert.Nil(t, err, "hmac verifying failed")
}
