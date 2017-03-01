package mock

import (
	"github.com/SermoDigital/jose/jws"
	"github.com/krinklesaurus/jwt_proxy"
)

func NewMockCore() app.CoreAuth {
	return &MockCore{}
}

type MockCore struct {
}

func (c *MockCore) PublicKey() (*app.PublicKey, error) {
	return &app.PublicKey{
		Keys: []string{"MOCK_PUBLIC_KEY"}}, nil
}

func (c *MockCore) Token(provider string, code string) (*app.Token, error) {
	return &app.Token{}, nil
}

func (c *MockCore) JwtToken(jws.Claims) ([]byte, error) {
	data := []byte("JWT_TOKEN")
	return data, nil
}

func (c *MockCore) Claims(token *app.Token) jws.Claims {
	return jws.Claims{}
}

func (c *MockCore) RedirectURI() string {
	return "REDIRECT_URI"
}

func (c *MockCore) AuthURL(provider string, state string) string {
	return "AUTH_URL"
}

func (c *MockCore) Providers() []string {
	return []string{"PROVIDER_1", "PROVIDER_2"}
}

func (c *MockCore) Auth(username string, password string) (*app.Token, error) {
	return &app.Token{}, nil
}

func (c *MockCore) LocalEnabled() bool {
	return true
}
