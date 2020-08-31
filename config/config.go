package config

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	app "github.com/krinklesaurus/jwt-proxy"
	"github.com/krinklesaurus/jwt-proxy/provider"
	"github.com/spf13/viper"
)

func Initialize(configFile string) (*app.Config, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var err error
	if err = viper.ReadInConfig(); err != nil {
		fmt.Println("could not read config file:", viper.ConfigFileUsed())
	}

	rootURI := viper.GetString("rootUri")
	if rootURI == "" {
		return nil, errors.New("no rootUri set")
	}

	redirectURI := viper.GetString("redirectUri")
	if redirectURI == "" {
		return nil, errors.New("no redirectUri set")
	}

	wwwRootDir := viper.GetString("wwwRootDir")
	if wwwRootDir == "" {
		wwwRootDir = "www"
	}

	providers := map[string]app.Provider{}

	googleConfig := viper.Sub("providers.google")
	if googleConfig != nil {
		providers["google"] = provider.NewGoogle(
			rootURI,
			viper.GetString("providers.google.clientId"),
			viper.GetString("providers.google.clientSecret"),
			viper.GetStringSlice("providers.google.scopes"),
		)
	}

	githubConfig := viper.Sub("providers.github")
	if githubConfig != nil {
		providers["github"] = provider.NewGithub(
			rootURI,
			viper.GetString("providers.github.clientId"),
			viper.GetString("providers.github.clientSecret"),
			viper.GetStringSlice("providers.github.scopes"),
		)
	}

	facebookConfig := viper.Sub("providers.facebook")
	if facebookConfig != nil {
		providers["facebook"] = provider.NewFacebook(
			rootURI,
			viper.GetString("providers.facebook.clientId"),
			viper.GetString("providers.facebook.clientSecret"),
			viper.GetStringSlice("providers.facebook.scopes"),
		)
	}

	if len(providers) <= 0 {
		return nil, errors.New("no providers have been configured")
	}

	audience := viper.GetString("jwt.audience")
	issuer := viper.GetString("jwt.issuer")
	subject := viper.GetString("jwt.subject")
	expiry := viper.GetInt("jwt.expirySeconds")

	signingMethodKey := viper.GetString("jwt.signingMethod")
	if signingMethodKey == "" {
		return nil, errors.New("no signing method set")
	}
	signingMethod := app.SigningMethods[signingMethodKey]
	if signingMethod == nil {
		return nil, errors.New("no valid signing method set")
	}

	var publicKeyData []byte
	publicKey := viper.GetString("jwt.publicKey")
	if publicKey == "" {
		publicKeyPath := viper.GetString("jwt.publicKeyPath")
		if publicKeyPath == "" {
			publicKeyPath = "certs/public.pem"
		}
		publicKeyData, err = ioutil.ReadFile(publicKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		publicKeyData = []byte(publicKey)
	}
	block, _ := pem.Decode(publicKeyData)
	rsaPub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var privateKeyData []byte
	privateKey := viper.GetString("jwt.privateKey")
	if privateKey == "" {
		privateKeyPath := viper.GetString("jwt.privateKeyPath")
		if privateKeyPath == "" {
			privateKeyPath = "certs/private.pem"
		}
		privateKeyData, err = ioutil.ReadFile(privateKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		privateKeyData = []byte(privateKey)
	}
	block2, _ := pem.Decode(privateKeyData)
	rsaPriv, err := x509.ParsePKCS1PrivateKey(block2.Bytes)
	if err != nil {
		return nil, err
	}

	return &app.Config{RootURI: rootURI,
		RedirectURI:   redirectURI,
		WWWRootDir:    wwwRootDir,
		Providers:     providers,
		SigningMethod: signingMethod,
		PrivateRSAKey: rsaPriv,
		PublicRSAKey:  rsaPub,
		Audience:      audience,
		Issuer:        issuer,
		Subject:       subject,
		ExpirySeconds: expiry}, nil
}
