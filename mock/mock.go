package mock

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/SermoDigital/jose/jwt"
)

var PrivateKey *rsa.PrivateKey
var PublicKey interface{}
var MockClaims jwt.Claims

func init() {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}

	PrivateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	pubKey := PrivateKey.PublicKey
	publicKeyDer, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		panic(err)
	}

	publicKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}

	PublicKey, err = x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	MockClaims = jwt.Claims{}
	MockClaims.SetIssuer("my-company")
	MockClaims.SetIssuedAt(time.Unix(1472828729, 0))
	MockClaims.SetExpiration(time.Unix(1504364729, 0))
	MockClaims.SetAudience("www.my-company.com")
	MockClaims.SetSubject("you@my-company.com")
}
