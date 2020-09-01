package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/krinklesaurus/jwt-proxy/log"
	"github.com/krinklesaurus/jwt-proxy/provider"
	"github.com/spf13/viper"
)

// Config is the struct for the config defined in a config.yml

type configRaw struct {
	RootURI     string `mapstructure:"rootUri"`
	RedirectURI string `mapstructure:"redirectUri"`
	WWWRootDir  string `mapstructure:"wwwRootDir"`
	Providers   map[string]struct {
		ClientID     string   `mapstructure:"clientId"`
		ClientSecret string   `mapstructure:"clientSecret"`
		Scopes       []string `mapstructure:"scopes"`
	} `mapstructure:"providers"`
	JWT struct {
		PrivateRSAKey     string `mapstructure:"privateRSAKey"`
		PublicRSAKey      string `mapstructure:"publicRSAKey"`
		PrivateRSAKeyPath string `mapstructure:"privateRSAKeyPath"`
		PublicRSAKeyPath  string `mapstructure:"publicRSAKeyPath"`
		SigningMethod     string `mapstructure:"signingMethod"`
		Audience          string `mapstructure:"audience"`
		Issuer            string `mapstructure:"issuer"`
		Subject           string `mapstructure:"subject"`
		ExpirySeconds     int    `mapstructure:"expirySeconds"`
	} `mapstructure:"jwt"`
}

type Config struct {
	RootURI           string
	RedirectURI       string
	WWWRootDir        string
	Providers         map[string]provider.Provider
	SigningMethod     string
	PrivateRSAKey     *rsa.PrivateKey
	PublicRSAKey      interface{}
	PrivateRSAKeyPath string
	PublicRSAKeyPath  string
	Audience          string
	Issuer            string
	Subject           string
	Password          string
	ExpirySeconds     int
}

func Initialize(configFile string) (*Config, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var err error
	if err = viper.ReadInConfig(); err != nil {
		fmt.Println("could not read config file:", viper.ConfigFileUsed())
	}

	var raw configRaw
	err = viper.Unmarshal(&raw)
	if err != nil {
		return nil, err
	}
	log.Debugf("loaded raw config: %+v", raw)

	providers := map[string]provider.Provider{}
	for key, providerInfo := range raw.Providers {
		switch key {
		case "google":
			providers["google"] = provider.NewGoogle(
				raw.RootURI,
				providerInfo.ClientID,
				providerInfo.ClientSecret,
				// bug in Viper/mapstructure forces us to split it manually again
				strings.Split(providerInfo.Scopes[0], " "),
			)
		case "facebook":
			providers["facebook"] = provider.NewFacebook(
				raw.RootURI,
				providerInfo.ClientID,
				providerInfo.ClientSecret,
				strings.Split(providerInfo.Scopes[0], " "),
			)
		case "github":
			providers["github"] = provider.NewGithub(
				raw.RootURI,
				providerInfo.ClientID,
				providerInfo.ClientSecret,
				strings.Split(providerInfo.Scopes[0], " "),
			)
		default:
			log.Warnf("no provider info for key %s", key)
		}
	}

	var publicKeyData []byte
	publicRSAKey := raw.JWT.PublicRSAKey
	publicRSAKeyPath := raw.JWT.PublicRSAKeyPath
	if publicRSAKey == "" {
		if publicRSAKeyPath == "" {
			publicRSAKeyPath = "certs/public.pem"
		}
		publicKeyData, err = ioutil.ReadFile(publicRSAKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		publicKeyData = []byte(publicRSAKey)
	}
	block, _ := pem.Decode(publicKeyData)
	if block == nil {
		return nil, fmt.Errorf("could not decode publicKeyData %v", string(publicKeyData))
	}
	rsaPub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var privateKeyData []byte
	privateRSAKey := raw.JWT.PrivateRSAKey
	privateRSAKeyPath := raw.JWT.PrivateRSAKeyPath
	if privateRSAKey == "" {
		if privateRSAKeyPath == "" {
			privateRSAKeyPath = "certs/private.pem"
		}
		privateKeyData, err = ioutil.ReadFile(privateRSAKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		privateKeyData = []byte(privateRSAKey)
	}
	block2, _ := pem.Decode(privateKeyData)
	rsaPriv, err := x509.ParsePKCS1PrivateKey(block2.Bytes)
	if err != nil {
		return nil, err
	}

	return &Config{RootURI: raw.RootURI,
		RedirectURI:   raw.RedirectURI,
		WWWRootDir:    raw.WWWRootDir,
		Providers:     providers,
		SigningMethod: raw.JWT.SigningMethod,
		PrivateRSAKey: rsaPriv,
		PublicRSAKey:  rsaPub,
		Audience:      raw.JWT.Audience,
		Issuer:        raw.JWT.Issuer,
		Subject:       raw.JWT.Subject,
		ExpirySeconds: raw.JWT.ExpirySeconds}, nil
}

// String is a helping toString function for the config for debugging
func (c *Config) String() string {
	providersString := ""
	for _, p := range c.Providers {
		providersString = providersString + fmt.Sprintf("%s with clientId %s, ", p.Name(), p.ClientID())
	}
	return fmt.Sprintf("rootURI: %s, redirectURI: %s, Audience: %s, Issuer: %s, Subject: %s, Providers: %s", c.RootURI, c.RedirectURI, c.Audience, c.Issuer, c.Subject, providersString)
}
