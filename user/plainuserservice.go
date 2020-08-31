package user

import (
	"fmt"
	"net/url"
)

type PlainUserService struct {
}

func (us *PlainUserService) UniqueUser(provider string, providerUserID string) (string, error) {
	return fmt.Sprintf("%s:%s", url.QueryEscape(provider), url.QueryEscape(providerUserID)), nil
}
