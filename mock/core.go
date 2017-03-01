package mock

import (
	"github.com/SermoDigital/jose/jws"
	"github.com/krinklesaurus/jwt_proxy"
)

func NewCore() app.CoreAuth {
	return &core{}
}

type core struct {
}

func (c *core) PublicKey() (*app.PublicKey, error) {
	return &app.PublicKey{
		Keys: []string{"MOCK_PUBLIC_KEY"}}, nil
}

func (c *core) Token(provider string, code string) (*app.Token, error) {
	return &app.Token{}, nil
}

func (c *core) JwtToken(jws.Claims) ([]byte, error) {
	data := []byte("JWT_TOKEN")
	return data, nil
}

func (c *core) Claims(token *app.Token) jws.Claims {
	return jws.Claims{}
}

func (c *core) RedirectURI() string {
	return "REDIRECT_URI"
}

func (c *core) AuthURL(provider string, state string) string {
	return "AUTH_URL"
}

func (c *core) Providers() []string {
	return []string{"PROVIDER_1", "PROVIDER_2"}
}

func (c *core) Auth(username string, password string) (*app.Token, error) {
	return &app.Token{}, nil
}

func (c *core) LocalEnabled() bool {
	return true
}
