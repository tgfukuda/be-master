package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	db "github.com/tgfukuda/be-master/db/sqlc"
	"github.com/tgfukuda/be-master/mocks"
	"github.com/tgfukuda/be-master/token"
	"github.com/tgfukuda/be-master/util"
)

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return mock.MatchedBy(func(given db.CreateUserParams) bool {
		err := util.CheckPassword(password, given.HashedPassword)
		if err != nil {
			return false
		}

		arg.HashedPassword = given.HashedPassword

		return reflect.DeepEqual(arg, given)
	})
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   "/users",
			method: http.MethodPost,
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					CreateUser(mock.Anything, EqCreateUserParams(db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					}, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusOK, recoder.Code)
				requireMatchUser(t, recoder.Body, user)
			},
		},
		{
			name:   "InvalidUserName",
			path:   "/users",
			method: http.MethodPost,
			body: gin.H{
				"username":  "abc#", // alphanum
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "PasswordTooShort",
			path:   "/users",
			method: http.MethodPost,
			body: gin.H{
				"username":  user.Username,
				"password":  "abc", // min: 6
				"full_name": user.FullName,
				"email":     user.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "EmptyFullName",
			path:   "/users",
			method: http.MethodPost,
			body: gin.H{
				"username":  "",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "InvalidPassword",
			path:   "/users",
			method: http.MethodPost,
			body: gin.H{
				"username":  user.Username,
				"password":  strings.Repeat("a", 73), // bcrypt default password length must be less than 73
				"full_name": user.FullName,
				"email":     user.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name:   "InvalidEmail",
			path:   "/users",
			method: http.MethodPost,
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "InternalError",
			path:   "/users",
			method: http.MethodPost,
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					CreateUser(mock.Anything, EqCreateUserParams(db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					}, password)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		// Add Forbidden test
	})
}

func TestLoginUser(t *testing.T) {
	user, password := randomUser(t)

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   "/users/login",
			method: http.MethodPost,
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					GetUser(mock.Anything, user.Username).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateNewSession(mock.Anything, mock.Anything).
					Times(1).
					Return(db.Session{}, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusOK, recoder.Code)
				requireMatchLoginUserResponse(t, tokenMaker, recoder.Body, user, password)
			},
		},
		{
			name:   "BadRequest",
			path:   "/users/login",
			method: http.MethodPost,
			body: gin.H{
				"username": "###",
				"password": password,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "NotFound",
			path:   "/users/login",
			method: http.MethodPost,
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					GetUser(mock.Anything, user.Username).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},
		{
			name:   "InternalError",
			path:   "/users/login",
			method: http.MethodPost,
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					GetUser(mock.Anything, user.Username).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name:   "UnAuthorizedPasswordUnMatch",
			path:   "/users/login",
			method: http.MethodPost,
			body: gin.H{
				"username": user.Username,
				"password": util.RandomString(6),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					GetUser(mock.Anything, user.Username).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
	})
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	assert.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(10),
		Email:          util.RandomEmail(),
	}

	return
}

func requireMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	var gotUser userResponse
	err = json.Unmarshal(data, &gotUser)
	assert.NoError(t, err)

	rsp := newUserResponse(user)
	assert.Equal(t, gotUser, rsp)
}

func requireMatchLoginUserResponse(t *testing.T, tokenMaker token.Maker, body *bytes.Buffer, user db.User, password string) {
	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	var gotRes loginUserResponse
	err = json.Unmarshal(data, &gotRes)
	assert.NoError(t, err)

	rsp := newUserResponse(user)

	assert.Equal(t, rsp, gotRes.User)
	payload, err := tokenMaker.VerifyToken(gotRes.AccessToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)
}
