package user

import "testing"

const hashedUserID string = "29e4c1b25d94c0379dab71eb2138f2a3c5171bd0fbfbc9f1ab2d94afe26afc95"

func TestCorrectHashedUserID(t *testing.T) {
	us := &HashUserService{}

	user, err := us.LoadUser("someProvider", "someProviderUserId")
	if err != nil {
		t.Error(err)
	}

	if user.ID != hashedUserID {
		t.Errorf("user.ID %s and hashedUserID %s do not match", user.ID, hashedUserID)
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
