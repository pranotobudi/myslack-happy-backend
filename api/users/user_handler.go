package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-happy-backend/common"
	"github.com/pranotobudi/myslack-happy-backend/mongodb"
)

type IUserHandler interface {
	GetUserByEmail(c *gin.Context)
	UserAuth(c *gin.Context)
	UpdateUserRooms(c *gin.Context)
	HelloWorld(c *gin.Context)
}
type userHandler struct {
	userService IUserService
}

// NewUserHandler will initialize userHandler object
func NewUserHandler() *userHandler {
	userService := NewUserService()
	return &userHandler{userService: userService}
}

// GetUserByEmail will return user based on email
func (h *userHandler) GetUserByEmail(c *gin.Context) {
	email, ok := c.GetQuery("email")
	log.Println("GetUserByEmail - email: ", email)
	if !ok {
		c.JSON(http.StatusBadRequest, email)
		return
	}
	// filter := bson.M{"email": email}
	// userPtr, err := mongo.GetUser(filter)
	userPtr, err := h.userService.GetUser(email)
	if err != nil {
		response := common.ResponseErrorFormatter(err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("inside room_io_handler-getRoom GetUserByEmail!: ", *userPtr)
	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
	log.Println("RESPONSE TO BROWSER: ", response)
	// Add CORS headers, if no global CORS setting
	// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	c.JSON(http.StatusOK, response)
}

// UserAuth will return user if exist or create new user if not exist
func (h *userHandler) UserAuth(c *gin.Context) {
	// login
	var userAuth mongodb.UserAuth

	err := c.BindJSON(&userAuth)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	log.Println("GetUserByEmail - email: ", userAuth.Email)
	userPtr, err := h.userService.UserAuth(userAuth)
	if err != nil {
		response := common.ResponseErrorFormatter(err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	fmt.Println("inside room_io_handler-UserAuth user registered! ID: ", *userPtr)
	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
	log.Println("RESPONSE TO BROWSER: ", response)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UpdateUserRooms(c *gin.Context) {
	// login
	var userMongo mongodb.User

	err := c.BindJSON(&userMongo)
	log.Println("UpdateUserRooms userMongo: ", userMongo)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userPtr, err := h.userService.UpdateUserRooms(userMongo)
	if err != nil {
		response := common.ResponseErrorFormatter(err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
	log.Println("RESPONSE TO BROWSER: ", response)
	c.JSON(http.StatusOK, response)

}

// HelloWorld will return welcome message for home path
func HelloWorld(c *gin.Context) {
	// firstname := c.DefaultQuery("firstname", "Guest")
	// lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

	// c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", "Hello from MySlack App. A Persistence Chat App... The server is running at the background..")
	log.Println("RESPONSE TO BROWSER: ", response)
	c.JSON(http.StatusOK, response)
}
