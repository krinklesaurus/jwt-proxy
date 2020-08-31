package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectPlainUserID(t *testing.T) {
	us := &PlainUserService{}

	userID, err := us.UniqueUser("someProvider", "someProviderUserId")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "someProvider:someProviderUserId", userID)

	userID, err = us.UniqueUser("somePro:vider", "someProviderUserId")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "somePro%3Avider:someProviderUserId", userID)
}
