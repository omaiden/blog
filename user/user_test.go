package user_test

import (
	"blog/pkg/tu"
	"blog/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister_Success(t *testing.T) {
	tc := tu.Setup()
	defer tc.Teardown()

	req := &user.RegisterRequest{
		Username:     "thunder",
		Email:        "thunder@noemail.com",
		PasswordHash: "pwd",
	}

	registeredUser, err := user.Register(tc.Ctx(), req)
	assert.NoError(t, err)
	assert.Equal(t, "thunder", registeredUser.Username)
	assert.Equal(t, "thunder@noemail.com", registeredUser.Email)
}

func TestRegister_InvalidRequest(t *testing.T) {
	tc := tu.Setup()
	defer tc.Teardown()

	req := &user.RegisterRequest{}

	_, err := user.Register(tc.Ctx(), req)
	assert.ErrorContains(t, err, "username required")
	assert.ErrorContains(t, err, "email required")
	assert.ErrorContains(t, err, "password required")
}

func TestRegister_UsernameExist(t *testing.T) {
	tc := tu.Setup()
	defer tc.Teardown()

	req := &user.RegisterRequest{
		Username:     "thunder",
		Email:        "thunder@noemail.com",
		PasswordHash: "pwd",
	}

	_, err := user.Register(tc.Ctx(), req)
	assert.NoError(t, err)

	_, err = user.Register(tc.Ctx(), req)
	assert.ErrorContains(t, err, "username already exists")
}
