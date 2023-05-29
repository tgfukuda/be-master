package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tgfukuda/be-master/util"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	assert.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)

	assert.Equal(t, payload.Username, username)
	assert.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	assert.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	assert.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrExpiredToken.Error())
	assert.Nil(t, payload)
}
