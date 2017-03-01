package mock

import (
	"fmt"

	"github.com/krinklesaurus/jwt_proxy"
)

func NewUserservice() app.UserService {
	return &userService{}
}

type userService struct {
}

func (us *userService) UniqueUser(provider string, providerUserID string) (app.User, error) {
	userID := fmt.Sprintf("%s", providerUserID)
	return app.User{
		ID: userID,
	}, nil
}

func (us *userService) LoginUser(username string, plainPassword string) error {
	return nil
}
