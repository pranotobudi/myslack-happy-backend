package messages

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-happy-backend/common"
	"go.mongodb.org/mongo-driver/bson"
)

type IMessageHandler interface {
	GetMessages(c *gin.Context)
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
func (h *messageHandler) GetMessages(c *gin.Context) {
	// TODO: we actually should use params instead of query, need changes in the frontend also
	// because room_id is unique resource (params), not filtering (query)
	// roomId := c.Param("room_id")
	// log.Println("GetMessages params: ", c.Params, "roomId: ", roomId)
	roomId, ok := c.GetQuery("room_id")
	log.Println("GetMessages - roomId: ", roomId, "Ok: ", ok)
	if !ok {
		c.JSON(http.StatusBadRequest, roomId)
		return
	}
	filter := bson.M{"room_id": roomId}
	// messages, err := mongo.GetMessages(filter)
	messages, err := h.service.GetMessages(filter)
	if err != nil {
		response := common.ResponseErrorFormatter(err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("inside room_io_handler-getMessages!: ", messages)
	response := common.ResponseFormatter(http.StatusOK, "success", "get messages successfull", messages)
	log.Println("RESPONSE TO BROWSER: ", response)
	// Add CORS headers, if no global CORS setting
	// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	c.JSON(http.StatusOK, response)
}
