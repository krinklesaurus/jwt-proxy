package mock

import (
	"fmt"

	app "github.com/krinklesaurus/jwt-proxy"
)

func NewUserservice() app.UserService {
	return &userService{}
}

type userService struct {
}

func (us *userService) UniqueUser(provider string, providerUserID string) (string, error) {
	userID := fmt.Sprintf("%s", providerUserID)
	return userID, nil
}

func (us *userService) LoginUser(username string, plainPassword string) error {
	return nil
}
