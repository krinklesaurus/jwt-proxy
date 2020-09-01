package mock

import (
	"encoding/json"
	"fmt"
	"time"

	app "github.com/krinklesaurus/jwt-proxy"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func newProvider() app.Provider {
	return &provider{}
}

type provider struct {
}

func (m *provider) AuthCodeURL(state string) string {
	return "MOCK_AUTH_CODE"
}

func (m *provider) UserInfo() (map[string]interface{}, string, error) {
	return nil, "MOCK_USER_ID", nil
}

func (m *provider) ClientID() string {
	return "CLIENT_ID"
}

func (m *provider) User() (string, error) {
	return "MOCK_USER_ID", nil
}

func (m *provider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  "MOCK_ACCESS_TOKEN",
		TokenType:    "MOCK_TOKEN_TYPE",
		RefreshToken: "MOCK_REFRESH_TOKEN",
		Expiry:       time.Now(),
	}, nil
}

func NewMockConfig() *app.Config {
	return &app.Config{
		PublicRSAKey:  PrivateKey.PublicKey,
		PrivateRSAKey: PrivateKey,
		Providers: map[string]app.Provider{
			"MOCK_PROVIDER": newProvider(),
		},
	}
}

func (p *provider) Name() string {
	return "MOCK_PROVIDER"
}

func (p *provider) String() string {
	toString := struct {
		ClientID   string   `json:"client_id"`
		AuthURL    string   `json:"auth_url"`
		TokenURL   string   `json:"token_url"`
		RediectURL string   `json:"redirect_url"`
		Scopes     []string `json:"scopes"`
	}{
		"mock-client-id",
		"mock-auth-url",
		"mock-token-url",
		"mock-redirect-url",
		[]string{"mock-scope-1", "mock-scope-2"},
	}
	b, err := json.Marshal(toString)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return string(b)
}
