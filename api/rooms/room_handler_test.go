package rooms

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
	getRoomsFunc   func() ([]mongodb.Room, error)
	getAnyRoomFunc func() (*mongodb.Room, error)
	addRoomFunc    func(name string) (string, error)
)

type mockService struct{}

// type mockService struct {
// 	mock.Mock
// }

func (m *mockService) GetRooms() ([]mongodb.Room, error) {
	return getRoomsFunc()
}
func (m *mockService) GetAnyRoom() (*mongodb.Room, error) {
	return getAnyRoomFunc()
}
func (m *mockService) AddRoom(name string) (string, error) {
	return addRoomFunc(name)
}

func TestGetRooms(t *testing.T) {

	tt := []struct {
		Name     string
		mockFunc func() ([]mongodb.Room, error)
		CodeWant int
	}{
		{
			Name: "GetRooms Success",
			mockFunc: func() ([]mongodb.Room, error) {
				return []mongodb.Room{}, nil
			},
			CodeWant: http.StatusOK,
		},
		{
			Name: "GetRooms Failed",
			mockFunc: func() ([]mongodb.Room, error) {
				return nil, errors.New("get rooms failed")
			},
			CodeWant: http.StatusInternalServerError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getRoomsFunc = tc.mockFunc

			// messageHandler := NewMessageHandler(&mockService{})
			roomHandler := NewRoomHandler()
			roomHandler.roomService = &mockService{}
			rc := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rc)
			c.Request, _ = http.NewRequest(http.MethodGet, "", nil) // c.Params doesn't 			log.Println(c.Params, c.Request.RequestURI)
			roomHandler.GetRooms(c)

			assert.EqualValues(t, tc.CodeWant, rc.Code)
			// log.Println("test response: ", rc.Body.String())
		})
	}
}

func TestGetAnyRoom(t *testing.T) {

	tt := []struct {
		Name     string
		mockFunc func() (*mongodb.Room, error)
		CodeWant int
	}{
		{
			Name: "GetAnyRoom Success",
			mockFunc: func() (*mongodb.Room, error) {
				return &mongodb.Room{}, nil
			},
			CodeWant: http.StatusOK,
		},
		{
			Name: "GetAnyRoom Failed",
			mockFunc: func() (*mongodb.Room, error) {
				return nil, errors.New("get any room failed")
			},
			CodeWant: http.StatusInternalServerError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getAnyRoomFunc = tc.mockFunc

			// messageHandler := NewMessageHandler(&mockService{})
			roomHandler := NewRoomHandler()
			roomHandler.roomService = &mockService{}
			rc := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rc)
			c.Request, _ = http.NewRequest(http.MethodGet, "", nil) // c.Params doesn't 			log.Println(c.Params, c.Request.RequestURI)
			roomHandler.GetAnyRoom(c)

			assert.EqualValues(t, tc.CodeWant, rc.Code)
			// log.Println("test response: ", rc.Body.String())
		})
	}
}

func TestAddRoom(t *testing.T) {

	tt := []struct {
		Name       string
		mockFunc   func(name string) (string, error)
		CodeWant   int
		HttpMethod string
		Body       []byte
	}{
		{
			Name: "AddRoom Success",
			mockFunc: func(name string) (string, error) {
				return "room1", nil
			},
			CodeWant:   http.StatusOK,
			HttpMethod: http.MethodPost,
			Body:       []byte(`{"id":"1", "name":"budi"}`),
		},
		{
			Name: "AddRoom Failed json format error",
			mockFunc: func(name string) (string, error) {
				return "room1", nil
			},
			CodeWant:   http.StatusBadRequest,
			HttpMethod: http.MethodPost,
			Body:       []byte(``),
		},
		{
			Name: "AddRoom Failed",
			mockFunc: func(name string) (string, error) {
				return "", errors.New("add room failed")
			},
			CodeWant:   http.StatusInternalServerError,
			HttpMethod: http.MethodPost,
			Body:       []byte(`{"id":"1", "name":"budi"}`),
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			addRoomFunc = tc.mockFunc

			// messageHandler := NewMessageHandler(&mockService{})
			roomHandler := NewRoomHandler()
			roomHandler.roomService = &mockService{}
			rc := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rc)
			c.Request, _ = http.NewRequest(tc.HttpMethod, "", bytes.NewBuffer(tc.Body)) // c.Params doesn't 			log.Println(c.Params, c.Request.RequestURI)
			roomHandler.AddRoom(c)

			assert.EqualValues(t, tc.CodeWant, rc.Code)
			// log.Println("test response: ", rc.Body.String())
		})
	}
}
