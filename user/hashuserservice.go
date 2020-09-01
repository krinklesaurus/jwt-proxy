package user

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/krinklesaurus/jwt-proxy/log"
)

type HashUserService struct {
}

func (us HashUserService) UniqueUser(provider string, providerUserID string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(provider + ":" + providerUserID))
	hashedUserID := fmt.Sprintf("%x", hash.Sum(nil))

	log.Debugf("user id %s from provider \"%s\" and id \"%s\"", hashedUserID, provider, providerUserID)

	return hashedUserID, nil
}

func (us HashUserService) LoginUser(username string, plainPassword string) error {
	testPassword, err := bcrypt.GenerateFromPassword([]byte("tester"), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(testPassword, []byte(plainPassword))
	if err != nil {
		return err
	}

	return nil
}
