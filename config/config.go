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

func readString(key string, def string) (string, error) {
	val := viper.GetString(key)
	if val == "" {
		if def != "" {
			return def, nil
		}
		return "", fmt.Errorf("config %s must not be empty", key)
	}
	return val, nil
}

func readInt(key string, def int) (int, error) {
	val := viper.GetInt(key)
	if val == 0 {
		if def != 0 {
			return def, nil
		}
		return 0, fmt.Errorf("config %s must not 0", key)
	}
	return val, nil
}

func readSlice(key string) ([]string, error) {
	val := viper.GetStringSlice(key)
	if len(val) == 0 {
		return []string{}, fmt.Errorf("config %s must not empty", key)
	}
	return val, nil
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

	rootURI, err := readString("rootUri", "")
	if err != nil {
		return nil, err
	}
	redirectURI, err := readString("redirectUri", "")
	if err != nil {
		return nil, err
	}
	wwwRootDir, err := readString("wwwRootDir", "www")
	if err != nil {
		return nil, err
	}

	publicRSAKey := viper.GetString("jwt.publicRSAKey")
	publicRSAKeyPath, err := readString("jwt.publicRSAKeyPath", "certs/public.pem")
	if err != nil {
		return nil, err
	}
	privateRSAKey := viper.GetString("jwt.privateRSAKey")
	privateRSAKeyPath, err := readString("jwt.privateRSAKeyPath", "certs/private.pem")
	if err != nil {
		return nil, err
	}
	signingMethod := viper.GetString("jwt.signingMethod")

	audience, err := readString("jwt.audience", "")
	if err != nil {
		return nil, err
	}
	subject, err := readString("jwt.subject", "")
	if err != nil {
		return nil, err
	}
	issuer, err := readString("jwt.issuer", "")
	if err != nil {
		return nil, err
	}
	expirySeconds, err := readInt("jwt.expirySeconds", 86400)
	if err != nil {
		return nil, err
	}

	providers := map[string]provider.Provider{}

	googleClientID := viper.GetString("providers.google.clientId")
	if googleClientID != "" {
		log.Debugf("found provider google")
		googleSecret := viper.GetString("providers.google.clientSecret")
		if googleSecret == "" {
			return nil, fmt.Errorf("google secret must not be empty")
		}
		googleScopes := viper.GetStringSlice("providers.google.scopes")
		if len(googleScopes) == 0 {
			return nil, fmt.Errorf("google scopes must not be empty")
		}
		providers["google"] = provider.NewGoogle(
			rootURI,
			googleClientID,
			googleSecret,
			googleScopes,
		)
	}

	githubClientID := viper.GetString("providers.github.clientId")
	if githubClientID != "" {
		log.Debugf("found provider github")
		githubSecret := viper.GetString("providers.github.clientSecret")
		if githubSecret == "" {
			return nil, fmt.Errorf("github secret must not be empty")
		}
		githubScopes := viper.GetStringSlice("providers.github.scopes")
		if len(githubScopes) == 0 {
			return nil, fmt.Errorf("github scopes must not be empty")
		}
		providers["github"] = provider.NewGithub(
			rootURI,
			githubClientID,
			githubSecret,
			githubScopes,
		)
	}

	facebookClientID := viper.GetString("providers.facebook.clientId")
	if facebookClientID != "" {
		log.Debugf("found provider facebook")
		facebookSecret := viper.GetString("providers.facebook.clientSecret")
		if facebookSecret == "" {
			return nil, fmt.Errorf("facebook secret must not be empty")
		}
		facebookScopes := viper.GetStringSlice("providers.facebook.scopes")
		if len(facebookScopes) == 0 {
			return nil, fmt.Errorf("facebook scopes must not be empty")
		}
		providers["facebook"] = provider.NewFacebook(
			rootURI,
			facebookClientID,
			viper.GetString("providers.facebook.clientSecret"),
			viper.GetStringSlice("providers.facebook.scopes"),
		)
	}

	if err != nil {
		return nil, err
	}

	var publicKeyData []byte
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

	return &Config{RootURI: rootURI,
		RedirectURI:       redirectURI,
		WWWRootDir:        wwwRootDir,
		Providers:         providers,
		SigningMethod:     signingMethod,
		PrivateRSAKey:     rsaPriv,
		PrivateRSAKeyPath: privateRSAKeyPath,
		PublicRSAKey:      rsaPub,
		PublicRSAKeyPath:  publicRSAKeyPath,
		Audience:          audience,
		Issuer:            issuer,
		Subject:           subject,
		ExpirySeconds:     expirySeconds}, nil
}

// String is a helping toString function for the config for debugging
func (c *Config) String() string {
	providersString := ""
	for _, p := range c.Providers {
		providersString = providersString + fmt.Sprintf("%s with clientId %s, ", p.Name(), p.ClientID())
	}
	return fmt.Sprintf("rootURI: %s, redirectURI: %s, WWWRootDir: %s, SigningMethod: %s, PublicRSAKeyPath: %s, PrivateKeyPath: %s, Audience: %s, Issuer: %s, Subject: %s, Expiry: %d, Providers: %s",
		c.RootURI, c.RedirectURI, c.WWWRootDir, c.SigningMethod, c.PublicRSAKeyPath, c.PrivateRSAKeyPath, c.Audience, c.Issuer, c.Subject, c.ExpirySeconds, providersString)
}
