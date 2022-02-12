package rooms

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/pranotobudi/myslack-happy-backend/common"
	"github.com/pranotobudi/myslack-happy-backend/mongodb"
)

type IRoomHandler interface {
	GetRooms(w http.ResponseWriter, r *http.Request)
	GetAnyRoom(w http.ResponseWriter, r *http.Request)
	AddRoom(w http.ResponseWriter, r *http.Request)
}

type roomHandler struct {
	roomService IRoomService
}

// NewRoomHandler will initialize roomHandler object
func NewRoomHandler() *roomHandler {
	roomService := NewRoomService()
	return &roomHandler{roomService: roomService}
}

// GetRooms will return all rooms available
func (h *roomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	// rooms, err := mongo.GetRooms()
	rooms, err := h.roomService.GetRooms()
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("inside room_io_handler-getRooms!: ", rooms)
	response := common.ResponseFormatter(http.StatusOK, "success", "get rooms successfull", rooms)
	log.Println("RESPONSE TO BROWSER: ", response)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	// w.Write([]byte(fmt.Sprintf("%v", response)))
	// c.JSON(http.StatusOK, response)
}

// GetAnyRoom will return one room with no specific condition
func (h *roomHandler) GetAnyRoom(w http.ResponseWriter, r *http.Request) {
	// request: userId
	// response: user snapshot to load main page
	roomPtr, err := h.roomService.GetAnyRoom()
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := common.ResponseFormatter(http.StatusOK, "success", "get rooms successfull", *roomPtr)
	log.Println("RESPONSE TO BROWSER: ", response)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	// w.Write([]byte(fmt.Sprintf("%v", response)))
	// c.JSON(http.StatusOK, response)
}

// AddRoom will add room to the database
func (h *roomHandler) AddRoom(w http.ResponseWriter, r *http.Request) {

	var room mongodb.Room
	// c.Bind(&roomName)
	err := json.NewDecoder(r.Body).Decode(&room)
	// err := c.BindJSON(&room)
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusBadRequest, errors.New("request Decoding failed"))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", room)))
		// c.JSON(http.StatusBadRequest, room)
		return
	}
	log.Println("JSON roomName: ", room.Name)
	// roomId, err := mongo.AddRoom(room.Name)
	roomId, err := h.roomService.AddRoom(room.Name)
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, errors.New("add room failed"))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", roomId)))
		// c.JSON(http.StatusInternalServerError, roomId)
		return
	}
	fmt.Println("room_io_handler-AddRoom: ", roomId)
	response := common.ResponseFormatter(http.StatusOK, "success", "add room successfull", roomId)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	// w.Write([]byte(fmt.Sprintf("%v", response)))
	// c.JSON(http.StatusOK, response)
}
