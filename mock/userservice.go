package mock

import (
	"fmt"

	"github.com/krinklesaurus/jwt_proxy"
)

func NewMockUserservice() app.UserService {
	return &MockUserService{}
}

type MockUserService struct {
}

func (us *MockUserService) LoadUser(provider string, providerUserID string) (app.User, error) {
	userID := fmt.Sprintf("%s", providerUserID)
	return app.User{
		ID: userID,
	}, nil
}

func (us *MockUserService) LoginUser(username string, plainPassword string) error {
	return nil
}
