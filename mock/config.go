package mock

import (
	"time"

	"github.com/krinklesaurus/jwt_proxy"
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

func (m *provider) UniqueUserID() (string, error) {
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
