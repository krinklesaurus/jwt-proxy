package config

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"

	"errors"

	"github.com/krinklesaurus/jwt_proxy"
	"github.com/krinklesaurus/jwt_proxy/provider"
	"github.com/spf13/viper"
)

const configName string = "config"
const configType string = "yaml"

const configRootURI string = "root_uri"
const configRedirectURI string = "redirect_uri"
const configProviders string = "providers"
const configProviderClientID string = "client_id"
const configProviderClientSecret string = "client_secret"
const configProviderScopes string = "scopes"

const configGoogleClientID string = configProviders + ".google." + configProviderClientID
const configGoogleClientSecret string = configProviders + ".google." + configProviderClientSecret
const configGoogleScopes string = configProviders + ".google." + configProviderScopes

const configFacebookClientID string = configProviders + ".facebook." + configProviderClientID
const configFacebookClientSecret string = configProviders + ".facebook." + configProviderClientSecret
const configFacebookScopes string = configProviders + ".facebook." + configProviderScopes

const configGithubClientID string = configProviders + ".github." + configProviderClientID
const configGithubClientSecret string = configProviders + ".github." + configProviderClientSecret
const configGithubScopes string = configProviders + ".github." + configProviderScopes

const configSigningMethod string = "jwt.signingMethod"
const configPublicKeyPath string = "jwt.public_key"
const configPrivateKeyPath string = "jwt.private_key"
const configJwtAudience string = "jwt.audience"
const configJwtIssuer string = "jwt.issuer"
const configJwtSubject string = "jwt.subject"

func Initialize(configFile string) (*app.Config, error) {
	viper.SetConfigType(configType)
	viper.AutomaticEnv()

	viper.BindEnv(configRootURI, "ROOT_URI")
	viper.BindEnv(configRedirectURI, "REDIRECT_URI")

	viper.BindEnv(configGoogleClientID, "GOOGLE_CLIENTID")
	viper.BindEnv(configGoogleClientSecret, "GOOGLE_SECRET")
	viper.BindEnv(configGoogleScopes, "GOOGLE_SCOPES")

	viper.BindEnv(configFacebookClientID, "FACEBOOK_CLIENTID")
	viper.BindEnv(configFacebookClientSecret, "FACEBOOK_SECRET")
	viper.BindEnv(configFacebookScopes, "FACEBOOK_SCOPES")

	viper.BindEnv(configGithubClientID, "GITHUB_CLIENTID")
	viper.BindEnv(configGithubClientSecret, "GITHUB_SECRET")
	viper.BindEnv(configGithubScopes, "GITHUB_SCOPES")

	viper.BindEnv(configSigningMethod, "SIGNINGMETHOD")
	viper.BindEnv(configPublicKeyPath, "PUBLICKEY_PATH")
	viper.BindEnv(configPrivateKeyPath, "PRIVATEKEY_PATH")

	viper.BindEnv(configJwtAudience, "JWT_AUDIENCE")
	viper.BindEnv(configJwtIssuer, "JWT_ISSUER")
	viper.BindEnv(configJwtSubject, "JWT_SUBJECT")

	configReader, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	err = viper.ReadConfig(configReader)
	if err != nil {
		return nil, err
	}

	rootURI := viper.GetString(configRootURI)
	if rootURI == "" {
		return nil, errors.New("No root_uri set!")
	}

	redirectURI := viper.GetString(configRedirectURI)
	if redirectURI == "" {
		return nil, errors.New("No redirect_uri set!")
	}

	providersConfig := viper.GetStringMap(configProviders)

	providers := map[string]app.Provider{}
	if providersConfig["google"] != "" {
		providers["google"] = provider.NewGoogle(
			rootURI,
			viper.GetString(configGoogleClientID),
			viper.GetString(configGoogleClientSecret),
			viper.GetStringSlice(configGoogleScopes),
		)
	}

	if providersConfig["github"] != "" {
		providers["github"] = provider.NewGithub(
			rootURI,
			viper.GetString(configGithubClientID),
			viper.GetString(configGithubClientSecret),
			viper.GetStringSlice(configGithubScopes),
		)
	}

	if providersConfig["facebook"] != "" {
		providers["facebook"] = provider.NewFacebook(
			rootURI,
			viper.GetString(configFacebookClientID),
			viper.GetString(configFacebookClientSecret),
			viper.GetStringSlice(configFacebookScopes),
		)
	}

	if len(providers) <= 0 {
		return nil, errors.New("No providers have been configured!")
	}

	audience := viper.GetString(configJwtAudience)
	issuer := viper.GetString(configJwtIssuer)
	subject := viper.GetString(configJwtSubject)

	signingMethodKey := viper.GetString(configSigningMethod)
	if signingMethodKey == "" {
		return nil, errors.New("No signing method set!")
	}
	signingMethod := app.SigningMethods[signingMethodKey]
	if signingMethod == nil {
		return nil, errors.New("No valid signing method set!")
	}

	publicKeyPath := viper.GetString(configPublicKeyPath)
	derBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(derBytes)
	rsaPub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKeyPath := viper.GetString(configPrivateKeyPath)
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
