package rooms

import (
	"errors"
	"testing"

	"github.com/pranotobudi/myslack-happy-backend/mongodb"
	"github.com/stretchr/testify/assert"
)

var (
	getRoomsRepoFunc   func() ([]mongodb.Room, error)
	getAnyRoomRepoFunc func() (*mongodb.Room, error)
	addRoomRepoFunc    func(name string) (string, error)
)

type mockRoomRepo struct {
	mongodb.IMongoDB
}

func (m *mockRoomRepo) GetRooms() ([]mongodb.Room, error) {
	return getRoomsRepoFunc()
}
func (m *mockRoomRepo) GetAnyRoom() (*mongodb.Room, error) {
	return getAnyRoomRepoFunc()
}
func (m *mockRoomRepo) AddRoom(name string) (string, error) {
	return addRoomRepoFunc(name)
}
func TestGetRoomsService(t *testing.T) {

	tt := []struct {
		Name      string
		mockFunc  func() ([]mongodb.Room, error)
		IsSuccess bool
	}{
		{
			Name: "GetRooms Success",
			mockFunc: func() ([]mongodb.Room, error) {
				return []mongodb.Room{}, nil
			},
			IsSuccess: true,
		},
		{
			Name: "GetRooms Failed",
			mockFunc: func() ([]mongodb.Room, error) {
				return nil, errors.New("Get rooms failed")
			},
			IsSuccess: false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getRoomsRepoFunc = tc.mockFunc
			roomService := NewRoomService()
			roomService.repo = &mockRoomRepo{}

			rooms, err := roomService.GetRooms()

			if tc.IsSuccess {
				assert.NotNil(t, rooms)
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, rooms)
			}
		})
	}

}

func TestGetAnyRoomService(t *testing.T) {

	tt := []struct {
		Name            string
		mockAnyRoomFunc func() (*mongodb.Room, error)
		IsSuccess       bool
	}{
		{
			Name: "GetAnyRoom Success",
			mockAnyRoomFunc: func() (*mongodb.Room, error) {
				return &mongodb.Room{ID: "61f61d94fc663b6f4c8f3172"}, nil
			},
			IsSuccess: true,
		},
		{
			Name: "GetAnyRoom Failed",
			mockAnyRoomFunc: func() (*mongodb.Room, error) {
				return nil, errors.New("Get Any room failed")
			},
			IsSuccess: false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getAnyRoomRepoFunc = tc.mockAnyRoomFunc
			roomService := NewRoomService()
			roomService.repo = &mockRoomRepo{}

			rooms, err := roomService.GetAnyRoom()

			if tc.IsSuccess {
				assert.NotNil(t, rooms)
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, rooms)
			}
		})
	}
}

func TestAddRoomService(t *testing.T) {

	tt := []struct {
		Name            string
		mockAddRoomFunc func(name string) (string, error)
		RoomName        string
		IsSuccess       bool
	}{
		{
			Name: "AddRoom Success",
			mockAddRoomFunc: func(name string) (string, error) {
				return "room1", nil
			},
			RoomName:  "room1",
			IsSuccess: true,
		},
		{
			Name: "AddRoom Failed",
			mockAddRoomFunc: func(name string) (string, error) {
				return "", errors.New("Failed to add room")
			},
			RoomName:  "room1",
			IsSuccess: false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			addRoomRepoFunc = tc.mockAddRoomFunc
			roomService := NewRoomService()
			roomService.repo = &mockRoomRepo{}

			room, err := roomService.AddRoom(tc.RoomName)

			if tc.IsSuccess {
				assert.NotNil(t, room)
				assert.Nil(t, err)
			} else {
				// assert.Equal(t, "")
				assert.NotNil(t, err)
			}
		})
	}
}
