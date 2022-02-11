package users

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-happy-backend/mongodb"
	"github.com/stretchr/testify/assert"
)

var (
	getUserFunc         func(email string) (*mongodb.User, error)
	userAuthFunc        func(userAuth mongodb.UserAuth) (*mongodb.User, error)
	updateUserRoomsFunc func(userMongo mongodb.User) (*mongodb.User, error)
)

type mockService struct{}

func (m *mockService) GetUser(email string) (*mongodb.User, error) {
	return getUserFunc(email)
}
func (m *mockService) UserAuth(userAuth mongodb.UserAuth) (*mongodb.User, error) {
	return userAuthFunc(userAuth)
}
func (m *mockService) UpdateUserRooms(userMongo mongodb.User) (*mongodb.User, error) {
	return updateUserRoomsFunc(userMongo)
}

func TestGetUser(t *testing.T) {

	tt := []struct {
		Name       string
		mockFunc   func(email string) (*mongodb.User, error)
		CodeWant   int
		HttpMethod string
	}{
		{
			Name: "GetUser Success",
			mockFunc: func(email string) (*mongodb.User, error) {
				return &mongodb.User{}, nil
			},
			CodeWant:   http.StatusOK,
			HttpMethod: http.MethodGet,
		},
		{
			Name: "GetUser Failed",
			mockFunc: func(email string) (*mongodb.User, error) {
				return nil, errors.New("get user failed")
			},
			CodeWant:   http.StatusInternalServerError,
			HttpMethod: http.MethodGet,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getUserFunc = tc.mockFunc

			// messageHandler := NewMessageHandler(&mockService{})
			userHandler := NewUserHandler()
			userHandler.userService = &mockService{}
			rc := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rc)
			c.Request, _ = http.NewRequest(tc.HttpMethod, "http://localhost:8080?email=bud@gmail.com", nil)

			userHandler.GetUserByEmail(c)

			assert.EqualValues(t, tc.CodeWant, rc.Code)
			// log.Println("test response: ", rc.Body.String())
		})
	}
}

func TestUserAuth(t *testing.T) {

	tt := []struct {
		Name       string
		mockFunc   func(userAuth mongodb.UserAuth) (*mongodb.User, error)
		CodeWant   int
		HttpMethod string
		Body       []byte
	}{
		{
			Name: "UserAuth Success",
			mockFunc: func(userAuth mongodb.UserAuth) (*mongodb.User, error) {
				return &mongodb.User{}, nil
			},
			CodeWant:   http.StatusOK,
			HttpMethod: http.MethodPost,
			Body:       []byte(`{"email":"bud@gmail.com", "user_image":"https://aws.com"}`),
		},
		{
			Name: "UserAuth Failed json format error",
			mockFunc: func(userAuth mongodb.UserAuth) (*mongodb.User, error) {
				return nil, errors.New("UserAuth Failed json format error")
			},
			CodeWant:   http.StatusBadRequest,
			HttpMethod: http.MethodPost,
			Body:       []byte(``),
		},
		{
			Name: "UserAuth Failed",
			mockFunc: func(userAuth mongodb.UserAuth) (*mongodb.User, error) {
				return nil, errors.New("UserAuth Failed")
			},
			CodeWant:   http.StatusInternalServerError,
			HttpMethod: http.MethodPost,
			Body:       []byte(`{"id":"1", "name":"budi"}`),
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			userAuthFunc = tc.mockFunc

			// messageHandler := NewMessageHandler(&mockService{})
			userHandler := NewUserHandler()
			userHandler.userService = &mockService{}
			rc := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rc)
			c.Request, _ = http.NewRequest(tc.HttpMethod, "", bytes.NewBuffer(tc.Body))

			userHandler.UserAuth(c)

			assert.EqualValues(t, tc.CodeWant, rc.Code)
			// log.Println("test response: ", rc.Body.String())
		})
	}
}

func TestUpdateUserRooms(t *testing.T) {

	tt := []struct {
		Name       string
		mockFunc   func(userMongo mongodb.User) (*mongodb.User, error)
		CodeWant   int
		HttpMethod string
		Body       []byte
	}{
		{
			Name: "UpdateUserRooms Success",
			mockFunc: func(userMongo mongodb.User) (*mongodb.User, error) {
				return &mongodb.User{}, nil
			},
			CodeWant:   http.StatusOK,
			HttpMethod: http.MethodPost,
			Body: []byte(`{
				"id":"61cfa908eca4dd2b9d11d9ee",
				"email": "bud@gmail.com",
				"username":"",
				"user_image":"https://lh3.googleusercontent.com/",
				"rooms":["61cc50877ea033031b1a950e"]
			}`),
		},
		{
			Name: "UpdateUserRooms Failed json format error",
			mockFunc: func(userMongo mongodb.User) (*mongodb.User, error) {
				return nil, errors.New("UpdateUserRooms Failed json format error")
			},
			CodeWant:   http.StatusBadRequest,
			HttpMethod: http.MethodPost,
			Body:       []byte(``),
		},
		{
			Name: "UpdateUserRooms Failed",
			mockFunc: func(userMongo mongodb.User) (*mongodb.User, error) {
				return nil, errors.New("UpdateUserRooms Failed")
			},
			CodeWant:   http.StatusInternalServerError,
			HttpMethod: http.MethodPost,
			Body: []byte(`{
				"id":"61cfa908eca4dd2b9d11d9ee",
				"email": "bud@gmail.com",
				"username":"",
				"user_image":"https://lh3.googleusercontent.com/",
				"rooms":["61cc50877ea033031b1a950e"]
			}`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			updateUserRoomsFunc = tc.mockFunc

			// messageHandler := NewMessageHandler(&mockService{})
			userHandler := NewUserHandler()
			userHandler.userService = &mockService{}
			rc := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rc)
			c.Request, _ = http.NewRequest(tc.HttpMethod, "", bytes.NewBuffer(tc.Body))

			userHandler.UpdateUserRooms(c)

			assert.EqualValues(t, tc.CodeWant, rc.Code)
			// log.Println("test response: ", rc.Body.String())
		})
	}
}
