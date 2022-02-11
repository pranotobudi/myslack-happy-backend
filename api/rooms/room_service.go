package rooms

import (
	"fmt"

	"github.com/pranotobudi/myslack-happy-backend/mongodb"
)

type IRoomService interface {
	GetRooms() ([]mongodb.Room, error)
	GetAnyRoom() (*mongodb.Room, error)
	AddRoom(name string) (string, error)
}
type roomService struct {
	repo mongodb.IMongoDB
}

// NewRoomService will initialize roomService object
func NewRoomService() *roomService {
	r := mongodb.NewMongoDB()
	return &roomService{repo: r}
}

// GetRooms will get all rooms available
func (s *roomService) GetRooms() ([]mongodb.Room, error) {
	return s.repo.GetRooms()
}

// GetAnyRoom will return one room with no specific condition
func (s *roomService) GetAnyRoom() (*mongodb.Room, error) {
	anyRoomPtr, err := s.repo.GetAnyRoom()
	if err != nil {
		return anyRoomPtr, err
	}
	fmt.Println("inside room_io_handler-getRoom anyRoom!: ", *anyRoomPtr)
	return anyRoomPtr, nil
	// objID, err := primitive.ObjectIDFromHex(anyRoomPtr.ID)
	// if err != nil {
	// 	// panic(err)
	// 	return anyRoomPtr, err
	// }

	// filter := bson.M{"_id": objID}
	// roomPtr, err := s.repo.GetRoom(filter)
	// if err != nil {
	// 	return roomPtr, err
	// }
	// fmt.Println("inside room_io_handler-getRoom!: ", *roomPtr)
	// return roomPtr, nil
}

// AddRoom will add room to the database
func (s *roomService) AddRoom(name string) (string, error) {
	return s.repo.AddRoom(name)
}
