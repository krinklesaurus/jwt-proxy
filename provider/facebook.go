package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/krinklesaurus/jwt-proxy/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

func NewFacebook(rootURI string, clientID string, clientSecret string, scopes []string) Provider {
	return &FacebookProvider{
		conf: oauth2.Config{
			RedirectURL:  rootURI + "/jwt-proxy/callback/facebook",
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			Endpoint:     facebook.Endpoint,
		},
		clientID: clientID,
	}
}

type FacebookProvider struct {
	conf     oauth2.Config
	token    *oauth2.Token
	clientID string
}

func (f *FacebookProvider) AuthCodeURL(state string) string {
	return f.conf.AuthCodeURL(state)
}

func (f *FacebookProvider) ClientID() string {
	return f.clientID
}

func (f *FacebookProvider) User() (string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v2.7/me?access_token=%s", f.token.AccessToken)

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	log.Debugf("contents from facebook: %s", contents)

	dec := json.NewDecoder(bytes.NewReader(contents))
	var asMap map[string]string
	dec.Decode(&asMap)
	return asMap["id"], nil
}

func (f *FacebookProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := f.conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	f.token = token
	return f.token, err
}

func (f *FacebookProvider) Name() string {
	return "facebook"
}

func (f *FacebookProvider) String() string {
	toString := struct {
		ClientID   string   `json:"client_id"`
		AuthURL    string   `json:"auth_url"`
		TokenURL   string   `json:"token_url"`
		RediectURL string   `json:"redirect_url"`
		Scopes     []string `json:"scopes"`
	}{
		f.conf.ClientID,
		f.conf.Endpoint.AuthURL,
		f.conf.Endpoint.TokenURL,
		f.conf.RedirectURL,
		f.conf.Scopes,
	}
	b, err := json.Marshal(toString)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return string(b)
}
