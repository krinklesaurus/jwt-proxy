package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/krinklesaurus/jwt_proxy"
	"github.com/krinklesaurus/jwt_proxy/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogle(rootURI string, clientID string, clientSecret string, scopes []string) app.Provider {
	return &GoogleProvider{conf: oauth2.Config{
		RedirectURL:  rootURI + "/callback/google",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}}
}

type GoogleProvider struct {
	conf  oauth2.Config
	token *oauth2.Token
}

func (g *GoogleProvider) AuthCodeURL(state string) string {
	return g.conf.AuthCodeURL(state)
}

func (g *GoogleProvider) UniqueUserID() (string, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", g.token.AccessToken)

	response, err := http.Get(url)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", nil
	}

	log.Debugf("contents from google: %s", contents)

	dec := json.NewDecoder(bytes.NewReader(contents))
	var asMap map[string]string
	dec.Decode(&asMap)
	return asMap["id"], nil
}

func (g *GoogleProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	g.token = token
	return g.token, err
}
