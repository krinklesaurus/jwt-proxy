package user

import (
	"fmt"
	"net/url"

	app "github.com/krinklesaurus/jwt_proxy"
)

type PlainUserService struct {
}

func (us *PlainUserService) UniqueUser(provider string, providerUserID string) (app.User, error) {
	return app.User{
		ID: fmt.Sprintf("%s:%s", url.QueryEscape(provider), url.QueryEscape(providerUserID)),
	}, nil
}
