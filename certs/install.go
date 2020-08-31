package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/krinklesaurus/jwt-proxy/config"
	"github.com/krinklesaurus/jwt-proxy/core"
	"github.com/krinklesaurus/jwt-proxy/log"
	"github.com/krinklesaurus/jwt-proxy/user"
	uuid "github.com/satori/go.uuid"
)

func installCerts() {
	f1, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	f2, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	pem.Encode(f1, &privateKeyBlock)

	publicKey := privateKey.PublicKey
	publicKeyDer, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}

	publicKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	pem.Encode(f2, &publicKeyBlock)

	f1.Close()
	f2.Close()
}

func createTestToken(c *core.Core) {
	claims := jws.Claims{}

	expiry := time.Now().AddDate(1, 0, 0)

	claims.SetAudience(c.Config.Audience)
	claims.SetExpiration(expiry)
	claims.SetIssuedAt(time.Now())
	claims.SetIssuer(c.Config.Issuer)
	claims.SetJWTID(uuid.NewV4().String())
	claims.SetNotBefore(time.Now())
	claims.SetSubject(c.Config.Subject)

	provider := "test"
	userID := uuid.NewV4().String()
	accessToken := uuid.NewV4().String()

	claims.Set("provider", provider)
	claims.Set("user", userID)
	claims.Set("access_token", accessToken)
	claims.Set("token_type", "Bearer")
	//claims.Set("refresh_token", nil)
	claims.Set("expiry", expiry.Unix())

	b, err := c.Tokenizer.Serialize(claims)
	if err != nil {
		fmt.Printf("Sorry, could not create test token: %v", err)
	}

	jwtAsString := string(b)

	log.Infof("--------- CREATE TEST TOKEN CLAIMS BEGIN ---------")
	log.Infof("%v", claims)
	log.Infof("--------- CREATED TEST TOKEN CLAIMS END ---------")

	log.Infof("--------- CREATED TEST TOKEN BEGIN ---------")
	log.Infof("%s", jwtAsString)
	log.Infof("--------- CREATED TEST TOKEN END ---------")
}

func main() {
	configPtr := flag.String("config", "", "configuration file")
	certs := flag.Bool("certs", false, "true if certs should be created from scratch")
	token := flag.Bool("token", false, "true if a test user token should be created")

	flag.Parse()

	if *certs {
		installCerts()
	}

	if *token {
		config, err := config.Initialize(*configPtr)
		if err != nil {
			panic(err)
		}
		log.Infof("Config initialized: %s", config.String())
		userService := &user.HashUserService{}
		tokenizer := core.NewRSATokenizer(config.SigningMethod, config.PrivateRSAKey)
		core := core.New(config, tokenizer, userService)
		createTestToken(core)
	}
}
