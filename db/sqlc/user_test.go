package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tgfukuda/be-master/util"
)

func createRandUser(t *testing.T) User {
	hp, err := util.HashPassword(util.RandomString(6))
	assert.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hp,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, arg.Username, user.Username)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.Equal(t, arg.FullName, user.FullName)
	assert.Equal(t, arg.Email, user.Email)

	assert.True(t, user.PasswordChangedAt.IsZero())
	assert.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	assert.NoError(t, err)
	assert.NotEmpty(t, user2)

	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, user1.HashedPassword, user2.HashedPassword)
	assert.Equal(t, user1.FullName, user2.FullName)
	assert.Equal(t, user1.Email, user2.Email)
	assert.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Millisecond)
	assert.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Millisecond)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	old := createRandUser(t)
	newFullName := util.RandomOwner()

	updated, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: old.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, updated)
	assert.Equal(t, old.Username, updated.Username)
	assert.Equal(t, old.HashedPassword, updated.HashedPassword)
	assert.Equal(t, newFullName, updated.FullName)
	assert.Equal(t, old.Email, updated.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	old := createRandUser(t)
	newEmail := util.RandomEmail()

	updated, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: old.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, updated)
	assert.Equal(t, old.Username, updated.Username)
	assert.Equal(t, old.HashedPassword, updated.HashedPassword)
	assert.Equal(t, old.FullName, updated.FullName)
	assert.Equal(t, newEmail, updated.Email)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	old := createRandUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	assert.NoError(t, err)

	updated, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: old.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, updated)
	assert.Equal(t, old.Username, updated.Username)
	assert.Equal(t, newHashedPassword, updated.HashedPassword)
	assert.Equal(t, old.FullName, updated.FullName)
	assert.Equal(t, old.Email, updated.Email)
}

func TestUpdateUserAll(t *testing.T) {
	old := createRandUser(t)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	assert.NoError(t, err)

	updated, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: old.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, updated)
	assert.Equal(t, old.Username, updated.Username)
	assert.Equal(t, newHashedPassword, updated.HashedPassword)
	assert.Equal(t, newFullName, updated.FullName)
	assert.Equal(t, newEmail, updated.Email)
}
