package messages

import (
	"github.com/pranotobudi/myslack-happy-backend/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type IMessageService interface {
	GetMessages(roomId string) ([]mongodb.Message, error)
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
func (s *messageService) GetMessages(roomId string) ([]mongodb.Message, error) {
	filter := bson.M{"room_id": roomId}
	messages, err := s.repo.GetMessages(filter)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
