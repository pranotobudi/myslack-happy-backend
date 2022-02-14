package messages

import (
	"errors"
	"testing"

	"github.com/pranotobudi/myslack-happy-backend/mongodb"
	"github.com/stretchr/testify/assert"
)

var (
	getMessagesRepoFunc func(filter interface{}) ([]mongodb.Message, error)
)

type mockMessageRepo struct {
	mongodb.IMongoDB
}

func (m *mockMessageRepo) GetMessages(filter interface{}) ([]mongodb.Message, error) {
	return getMessagesRepoFunc(filter)
}
func TestGetMessagesService(t *testing.T) {

	tt := []struct {
		Name      string
		mockFunc  func(filter interface{}) ([]mongodb.Message, error)
		roomId    string
		IsSuccess bool
	}{
		{
			Name: "GetMessages Success",
			mockFunc: func(filter interface{}) ([]mongodb.Message, error) {
				return []mongodb.Message{}, nil
			},
			roomId:    "abc1234567",
			IsSuccess: true,
		},
		{
			Name: "GetMessages Failed",
			mockFunc: func(filter interface{}) ([]mongodb.Message, error) {
				return nil, errors.New("get messages failed")
			},
			roomId:    "",
			IsSuccess: false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getMessagesRepoFunc = tc.mockFunc
			messageService := NewMessageService()
			messageService.repo = &mockMessageRepo{}

			messages, err := messageService.GetMessages(tc.roomId)

			if tc.IsSuccess {
				assert.NotNil(t, messages)
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, messages)
			}
		})
	}

}
