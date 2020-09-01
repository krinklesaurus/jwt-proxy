package core

import (
	"crypto/rsa"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
)

func NewRSATokenizer(signingMethod crypto.SigningMethod, privKey *rsa.PrivateKey) Tokenizer {
	return &RSATokenizer{signingMethod: signingMethod, privKey: privKey}
}

type RSATokenizer struct {
	signingMethod crypto.SigningMethod
	privKey       *rsa.PrivateKey
}

func (t *RSATokenizer) Serialize(claims map[string]interface{}) ([]byte, error) {
	jwt := jws.NewJWT(claims, t.signingMethod)
	return jwt.Serialize(t.privKey)
}
