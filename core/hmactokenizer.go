package core

import (
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"

	"github.com/krinklesaurus/jwt_proxy"
)

func NewHMACTokenizer(signingMethod crypto.SigningMethod, hmacKey []byte) app.Tokenizer {
	return &HMACTokenizer{signingMethod: signingMethod, hmacKey: hmacKey}
}

type HMACTokenizer struct {
	signingMethod crypto.SigningMethod
	hmacKey       []byte
}

func (t *HMACTokenizer) Serialize(claims map[string]interface{}) ([]byte, error) {
	jwt := jws.NewJWT(claims, t.signingMethod)
	return jwt.Serialize(t.hmacKey)
}
