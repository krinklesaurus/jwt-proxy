package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	app "github.com/krinklesaurus/jwt_proxy"
	"github.com/krinklesaurus/jwt_proxy/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func NewGithub(rootURI string, clientID string, clientSecret string, scopes []string) app.Provider {
	return &GithubProvider{conf: oauth2.Config{
		RedirectURL:  rootURI + "/callback/github",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint:     github.Endpoint,
	}}
}

type GithubProvider struct {
	conf  oauth2.Config
	token *oauth2.Token
}

func (g *GithubProvider) AuthCodeURL(state string) string {
	return g.conf.AuthCodeURL(state)
}

func (g *GithubProvider) UniqueUserID() (string, error) {
	url := fmt.Sprintf("https://api.github.com/user?access_token=%s", g.token.AccessToken)

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	log.Debugf("contents from github: %s", contents)

	dec := json.NewDecoder(bytes.NewReader(contents))
	var asMap map[string]string
	dec.Decode(&asMap)
	return asMap["login"], nil
}

func (g *GithubProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	g.token = token
	return g.token, err
}

func (g *GithubProvider) Name() string {
	return "github"
}

func (g *GithubProvider) String() string {
	toString := struct {
		ClientID   string   `json:"client_id"`
		AuthURL    string   `json:"auth_url"`
		TokenURL   string   `json:"token_url"`
		RediectURL string   `json:"redirect_url"`
		Scopes     []string `json:"scopes"`
	}{
		g.conf.ClientID,
		g.conf.Endpoint.AuthURL,
		g.conf.Endpoint.TokenURL,
		g.conf.RedirectURL,
		g.conf.Scopes,
	}
	b, err := json.Marshal(toString)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return string(b)
}
