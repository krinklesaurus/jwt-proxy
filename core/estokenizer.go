package core

import (
	"crypto/ecdsa"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/krinklesaurus/jwt_proxy"
)

func NewESTokenizer(signingMethod crypto.SigningMethod, privKey *ecdsa.PrivateKey) app.Tokenizer {
	return &ESTokenizer{signingMethod: signingMethod, privKey: privKey}
}

type ESTokenizer struct {
	signingMethod crypto.SigningMethod
	privKey       *ecdsa.PrivateKey
}

func (t *ESTokenizer) Serialize(claims map[string]interface{}) ([]byte, error) {
	jwt := jws.NewJWT(claims, t.signingMethod)
	return jwt.Serialize(t.privKey)
}
