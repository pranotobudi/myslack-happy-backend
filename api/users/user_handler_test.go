package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pranotobudi/myslack-happy-backend/common"
	"github.com/pranotobudi/myslack-happy-backend/mongodb"
	"github.com/stretchr/testify/assert"
)

var (
	getUserFunc         func(email string) (*mongodb.User, error)
	userAuthFunc        func(userAuth mongodb.UserAuth) (*mongodb.User, error)
	updateUserRoomsFunc func(userMongo mongodb.User) (*mongodb.User, error)
	userMailChatFunc    func(userMongo mongodb.User) (string, error)
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

func (m *mockService) UserMailChat(userMongo mongodb.User) (string, error) {
	return userMailChatFunc(userMongo)
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
			rr := httptest.NewRecorder()
			// c, _ := gin.CreateTestContext(rc)
			req, _ := http.NewRequest(tc.HttpMethod, "http://localhost:8080?email=bud@gmail.com", nil)

			userHandler.GetUserByEmail(rr, req)

			// check header StatusCode
			assert.EqualValues(t, tc.CodeWant, rr.Code)
			// check response (JSON format) StatusCode
			var response common.Response
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				assert.Errorf(t, err, "response format is not valid")
			}
			assert.EqualValues(t, tc.CodeWant, response.Meta.Code)
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
			rr := httptest.NewRecorder()
			// c, _ := gin.CreateTestContext(rc)
			req, _ := http.NewRequest(tc.HttpMethod, "", bytes.NewBuffer(tc.Body))

			userHandler.UserAuth(rr, req)

			// check header StatusCode
			assert.EqualValues(t, tc.CodeWant, rr.Code)
			// check response (JSON format) StatusCode
			var response common.Response
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				assert.Errorf(t, err, "response format is not valid")
			}
			assert.EqualValues(t, tc.CodeWant, response.Meta.Code)
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
			rr := httptest.NewRecorder()
			// c, _ := gin.CreateTestContext(rc)
			req, _ := http.NewRequest(tc.HttpMethod, "", bytes.NewBuffer(tc.Body))

			userHandler.UpdateUserRooms(rr, req)

			// check header StatusCode
			assert.EqualValues(t, tc.CodeWant, rr.Code)
			// check response (JSON format) StatusCode
			var response common.Response
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				assert.Errorf(t, err, "response format is not valid")
			}
			assert.EqualValues(t, tc.CodeWant, response.Meta.Code)
			// log.Println("test response: ", rc.Body.String())
		})
	}
}
