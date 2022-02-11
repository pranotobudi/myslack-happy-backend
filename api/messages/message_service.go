package messages

import (
	"github.com/pranotobudi/myslack-happy-backend/mongodb"
)

type IMessageService interface {
	GetMessages(filter interface{}) ([]mongodb.Message, error)
}
type messageService struct {
	repo mongodb.IMongoDB
}

// NewMessageService will initialize messageService object
func NewMessageService() *messageService {
	r := mongodb.NewMongoDB()
	return &messageService{repo: r}
}

// GetMessages will get messages based on the filter argument
func (s *messageService) GetMessages(filter interface{}) ([]mongodb.Message, error) {
	messages, err := s.repo.GetMessages(filter)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
