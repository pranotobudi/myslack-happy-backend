package messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/pranotobudi/myslack-happy-backend/common"
	"go.mongodb.org/mongo-driver/bson"
)

type IMessageHandler interface {
	GetMessages(w http.ResponseWriter, r *http.Request)
}
type messageHandler struct {
	service IMessageService
}

// NewMessageHandler initialize messageHandler object
func NewMessageHandler() *messageHandler {

	// func NewMessageHandler(messageService IMessageService) *messageHandler {
	messageService := NewMessageService()
	return &messageHandler{service: messageService}
}

// GetMessages will return list of messages for a room_id
func (h *messageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	// roomId := chi.URLParam(r, "room_id")
	roomId := r.URL.Query().Get("room_id")
	log.Println("GetMessages - roomId: ", roomId)
	if roomId == "" {
		response := common.ResponseErrorFormatter(http.StatusBadRequest, errors.New(roomId))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", roomId)))
		return
	}
	filter := bson.M{"room_id": roomId}
	messages, err := h.service.GetMessages(filter)
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		return
	}
	fmt.Println("inside room_io_handler-getMessages!: ", messages)
	response := common.ResponseFormatter(http.StatusOK, "success", "get messages successfull", messages)
	log.Println("RESPONSE TO BROWSER: ", response)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	// w.Write([]byte(fmt.Sprintf("%v", response)))
	// w.Write([]byte("kadlskfjal"))
}
