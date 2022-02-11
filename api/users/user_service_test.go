package users

import (
	"errors"
	"testing"

	"github.com/pranotobudi/myslack-happy-backend/mongodb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	getUserRepoFunc    func(filter interface{}) (*mongodb.User, error)
	addUserRepoFunc    func(user interface{}) (string, error)
	updateUserRepoFunc func(filter interface{}, update interface{}, options *options.UpdateOptions) error
)

type mockUserRepo struct {
	mongodb.IMongoDB
}

func (m *mockUserRepo) GetUser(filter interface{}) (*mongodb.User, error) {
	return getUserRepoFunc(filter)
}

func (m *mockUserRepo) AddUser(user interface{}) (string, error) {
	return addUserRepoFunc(user)
}
func (m *mockUserRepo) UpdateUser(filter interface{}, update interface{}, options *options.UpdateOptions) error {
	return updateUserRepoFunc(filter, update, options)
}
func TestGetUserService(t *testing.T) {

	tt := []struct {
		Name      string
		mockFunc  func(filter interface{}) (*mongodb.User, error)
		IsSuccess bool
	}{
		{
			Name: "GetUser Success",
			mockFunc: func(filter interface{}) (*mongodb.User, error) {
				return &mongodb.User{}, nil
			},
			IsSuccess: true,
		},
		{
			Name: "GetUser Failed",
			mockFunc: func(filter interface{}) (*mongodb.User, error) {
				return nil, errors.New("get user failed")
			},
			IsSuccess: false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getUserRepoFunc = tc.mockFunc
			userService := NewUserService()
			userService.repo = &mockUserRepo{}

			user, err := userService.GetUser("lumion@gmail.com")

			if tc.IsSuccess {
				assert.NotNil(t, user)
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, user)
			}
		})
	}
}

func TestUserAuthService(t *testing.T) {

	tt := []struct {
		Name            string
		getUserMockFunc func(filter interface{}) (*mongodb.User, error)
		addUserMockFunc func(user interface{}) (string, error)
		IsSuccess       bool
	}{
		{
			Name: "UserAuth Success",
			getUserMockFunc: func(filter interface{}) (*mongodb.User, error) {
				return &mongodb.User{ID: "abc123"}, nil
			},
			addUserMockFunc: func(user interface{}) (string, error) {
				return "abc123", nil
			},
			IsSuccess: true,
		},
		{
			Name: "UserAuth Failed",
			getUserMockFunc: func(filter interface{}) (*mongodb.User, error) {
				return nil, errors.New("get user failed")
			},
			addUserMockFunc: func(user interface{}) (string, error) {
				return "", errors.New("add user failed")
			},
			IsSuccess: false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getUserRepoFunc = tc.getUserMockFunc
			addUserRepoFunc = tc.addUserMockFunc
			userService := NewUserService()
			userService.repo = &mockUserRepo{}

			user, err := userService.UserAuth(mongodb.UserAuth{Email: "lumion@gmail.com"})

			if tc.IsSuccess {
				assert.NotNil(t, user)
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, user)
			}
		})
	}
}

func TestUpdateUserRoomsService(t *testing.T) {

	tt := []struct {
		Name               string
		getUserMockFunc    func(filter interface{}) (*mongodb.User, error)
		updateUserMockFunc func(filter interface{}, update interface{}, options *options.UpdateOptions) error
		IsSuccess          bool
	}{
		{
			Name: "UpdateUserRooms Success",
			getUserMockFunc: func(filter interface{}) (*mongodb.User, error) {
				return &mongodb.User{ID: "abc123"}, nil
			},
			updateUserMockFunc: func(filter interface{}, update interface{}, options *options.UpdateOptions) error {
				return nil
			},
			IsSuccess: true,
		},
		{
			Name: "UpdateUserRooms Failed",
			getUserMockFunc: func(filter interface{}) (*mongodb.User, error) {
				return nil, errors.New("get User is failed")
			},
			updateUserMockFunc: func(filter interface{}, update interface{}, options *options.UpdateOptions) error {
				return errors.New("update User is failed")
			},
			IsSuccess: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getUserRepoFunc = tc.getUserMockFunc
			updateUserRepoFunc = tc.updateUserMockFunc
			userService := NewUserService()
			userService.repo = &mockUserRepo{}

			user, err := userService.UpdateUserRooms(mongodb.User{ID: "abd123"})

			if tc.IsSuccess {
				assert.NotNil(t, user)
				assert.Nil(t, err)
			} else {
				assert.Nil(t, user)
				assert.NotNil(t, err)
			}
		})
	}
}
