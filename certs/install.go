package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/krinklesaurus/jwt-proxy/config"
	"github.com/krinklesaurus/jwt-proxy/core"
	"github.com/krinklesaurus/jwt-proxy/log"
	"github.com/krinklesaurus/jwt-proxy/user"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

func installCerts(config *config.Config) {
	f1, err := os.Create(config.PrivateRSAKeyPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("file %s could not be opened", config.PrivateRSAKeyPath))
		panic(err)
	}
	f2, err := os.Create(config.PublicRSAKeyPath)
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

func createTestToken(config *config.Config) {
	userService := &user.HashUserService{}
	tokenizer := core.NewRSATokenizer(core.SigningMethods[config.SigningMethod], config.PrivateRSAKey)

	provider := pseudoProvider{}
	oauthToken, _ := provider.Exchange(context.Background(), uuid.NewV4().String())
	user, _ := provider.User()
	c := core.New(config, tokenizer, userService)
	token := &core.TokenInfo{
		Token:    *oauthToken,
		User:     user,
		Provider: provider,
	}
	claims, _ := c.Claims(token)

	b, _ := c.JwtToken(claims)

	jwtAsString := string(b)

	log.Infof("--------- CREATE TEST TOKEN CLAIMS BEGIN ---------")
	log.Infof("%v", claims)
	log.Infof("--------- CREATED TEST TOKEN CLAIMS END ---------")

	log.Infof("--------- CREATED TEST TOKEN BEGIN ---------")
	log.Infof("%s", jwtAsString)
	log.Infof("--------- CREATED TEST TOKEN END ---------")
}

type pseudoProvider struct {
}

func (p pseudoProvider) AuthCodeURL(state string) string {
	return ""
}

func (p pseudoProvider) User() (string, error) {
	return fmt.Sprintf("user-%s", uuid.NewV4().String()), nil
}

func (p pseudoProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  uuid.NewV4().String(),
		TokenType:    "Bearer",
		RefreshToken: uuid.NewV4().String(),
		Expiry:       time.Now().Add(time.Duration(24*30) * time.Hour),
	}, nil
}

func (p pseudoProvider) String() string {
	return p.Name()
}

func (p pseudoProvider) Name() string {
	return "pseudo_provider"
}

func (p pseudoProvider) ClientID() string {
	return "pseudo_client_id"
}

func main() {
	configPtr := flag.String("config", "", "configuration file")
	certs := flag.Bool("certs", false, "true if certs should be created from scratch")
	token := flag.Bool("token", false, "true if a test user token should be created")

	flag.Parse()

	viper.SetConfigFile(*configPtr)
	var err error
	if err = viper.ReadInConfig(); err != nil {
		fmt.Println("could not read config file:", viper.ConfigFileUsed())
	}

	config := &config.Config{
		Audience:          viper.GetString("jwt.audience"),
		Issuer:            viper.GetString("jwt.issuer"),
		ExpirySeconds:     viper.GetInt("jwt.expirySeconds"),
		SigningMethod:     viper.GetString("jwt.signingMethod"),
		PrivateRSAKeyPath: viper.GetString("jwt.privateRSAKeyPath"),
		PublicRSAKeyPath:  viper.GetString("jwt.publicRSAKeyPath"),
	}

	if *certs {
		installCerts(config)
	}

	if *token {
		createTestToken(config)
	}
}
