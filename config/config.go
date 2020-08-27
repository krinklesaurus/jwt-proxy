package config

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	app "github.com/krinklesaurus/jwt_proxy"
	"github.com/krinklesaurus/jwt_proxy/provider"
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

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("could read config file:", viper.ConfigFileUsed())
	}

	rootURI := viper.GetString("root-uri")
	if rootURI == "" {
		return nil, errors.New("no root-uri set")
	}

	redirectURI := viper.GetString("redirect-uri")
	if redirectURI == "" {
		return nil, errors.New("no redirect-uri set")
	}

	providers := map[string]app.Provider{}

	googleConfig := viper.Sub("google")
	if googleConfig != nil {
		providers["google"] = provider.NewGoogle(
			rootURI,
			viper.GetString("google.clientid"),
			viper.GetString("google.clientsecret"),
			viper.GetStringSlice("google.scopes"),
		)
	}

	githubConfig := viper.Sub("github")
	if githubConfig != nil {
		providers["github"] = provider.NewGithub(
			rootURI,
			viper.GetString("github.clientid"),
			viper.GetString("github.clientsecret"),
			viper.GetStringSlice("github.scopes"),
		)
	}

	facebookConfig := viper.Sub("facebook")
	if facebookConfig != nil {
		providers["facebook"] = provider.NewFacebook(
			rootURI,
			viper.GetString("facebook.clientid"),
			viper.GetString("facebook.clientsecret"),
			viper.GetStringSlice("facebook.scopes"),
		)
	}

	if len(providers) <= 0 {
		return nil, errors.New("no providers have been configured")
	}

	audience := viper.GetString("jwt.audience")
	issuer := viper.GetString("jwt.issuer")
	subject := viper.GetString("jwt.subject")

	signingMethodKey := viper.GetString("jwt.signingMethod")
	if signingMethodKey == "" {
		return nil, errors.New("no signing method set")
	}
	signingMethod := app.SigningMethods[signingMethodKey]
	if signingMethod == nil {
		return nil, errors.New("no valid signing method set")
	}

	publicKeyPath := viper.GetString("jwt.public-key")
	if publicKeyPath == "" {
		publicKeyPath = "certs/public.pem"
	}
	derBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(derBytes)
	rsaPub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKeyPath := viper.GetString("jwt.private-key")
	if privateKeyPath == "" {
		privateKeyPath = "certs/private.pem"
	}
	der, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	block2, _ := pem.Decode(der)
	rsaPriv, err := x509.ParsePKCS1PrivateKey(block2.Bytes)
	if err != nil {
		return nil, err
	}

	return &app.Config{RootURI: rootURI,
		RedirectURI:   redirectURI,
		Providers:     providers,
		SigningMethod: signingMethod,
		PrivateRSAKey: rsaPriv,
		PublicRSAKey:  rsaPub,
		Audience:      audience,
		Issuer:        issuer,
		Subject:       subject}, nil
}
