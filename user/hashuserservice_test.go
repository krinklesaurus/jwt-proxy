package user

import "testing"

const hashedUserID string = "29e4c1b25d94c0379dab71eb2138f2a3c5171bd0fbfbc9f1ab2d94afe26afc95"

func TestCorrectHashedUserID(t *testing.T) {
	us := &HashUserService{}

	userID, err := us.UniqueUser("someProvider", "someProviderUserId")
	if err != nil {
		t.Error(err)
	}

	if userID != hashedUserID {
		t.Errorf("userID %s and hashedUserID %s do not match", userID, hashedUserID)
	}
}

func TestLoginTester(t *testing.T) {
	us := &HashUserService{}

	err := us.LoginUser("tester", "tester")
	if err != nil {
		t.Error(err)
	}

	err = us.LoginUser("tester", "wrong")
	if err == nil {
		t.Error("Could log in with wrong password \"wrong\"")
	}
}
