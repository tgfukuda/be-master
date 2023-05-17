package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(10)

	hashedPassword1, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword1)

	err = CheckPassword(password, hashedPassword1)
	assert.NoError(t, err)

	wrongPassword := RandomString(5)

	err = CheckPassword(wrongPassword, hashedPassword1)
	assert.Equal(t, err, bcrypt.ErrMismatchedHashAndPassword)

	hashedPassword2, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword2)
	assert.NotEqual(t, hashedPassword1, hashedPassword2)
}
