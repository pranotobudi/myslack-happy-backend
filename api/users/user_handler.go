package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/pranotobudi/myslack-happy-backend/api/rooms"
	"github.com/pranotobudi/myslack-happy-backend/common"
	"github.com/pranotobudi/myslack-happy-backend/mongodb"
)

type IUserHandler interface {
	GetUserByEmail(w http.ResponseWriter, r *http.Request)
	UserAuth(w http.ResponseWriter, r *http.Request)
	UpdateUserRooms(w http.ResponseWriter, r *http.Request)
	HelloWorld(w http.ResponseWriter, r *http.Request)
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
func (h *userHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	// email, ok := c.GetQuery("email")
	// email := chi.URLParam(r, "email")
	email := r.URL.Query().Get("email")
	log.Println("GetUserByEmail - email: ", email)

	if email == "" {
		response := common.ResponseErrorFormatter(http.StatusBadRequest, errors.New("failed to get email param"))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusBadRequest, email)
		return
	}
	// filter := bson.M{"email": email}
	// userPtr, err := mongo.GetUser(filter)
	userPtr, err := h.userService.GetUser(email)
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("inside room_io_handler-getRoom GetUserByEmail!: ", *userPtr)
	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
	log.Println("RESPONSE TO BROWSER: ", response)
	// Add CORS headers, if no global CORS setting
	// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	// w.Write([]byte(fmt.Sprintf("%v", response)))
	// c.JSON(http.StatusOK, response)
}

// UserAuth will return user if exist or create new user if not exist
func (h *userHandler) UserAuth(w http.ResponseWriter, r *http.Request) {
	// login
	var userAuth mongodb.UserAuth

	err := json.NewDecoder(r.Body).Decode(&userAuth)
	// err := c.BindJSON(&userAuth)

	if err != nil {
		// c.JSON(http.StatusBadRequest, err)
		response := common.ResponseErrorFormatter(http.StatusBadRequest, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		return
	}
	log.Println("GetUserByEmail - email: ", userAuth.Email)
	userPtr, err := h.userService.UserAuth(userAuth)
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusInternalServerError, response)
		return
	}

	fmt.Println("inside room_io_handler-UserAuth user registered! ID: ", *userPtr)
	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
	log.Println("RESPONSE TO BROWSER: ", response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	// w.Write([]byte(fmt.Sprintf("%v", response)))
	// c.JSON(http.StatusOK, response)
}

func (h *userHandler) UpdateUserRooms(w http.ResponseWriter, r *http.Request) {
	// login
	var userMongo mongodb.User
	err := json.NewDecoder(r.Body).Decode(&userMongo)

	// err := c.BindJSON(&userMongo)
	log.Println("UpdateUserRooms userMongo: ", userMongo)
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusBadRequest, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusBadRequest, err)
		return
	}

	userPtr, err := h.userService.UpdateUserRooms(userMongo)
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
	log.Println("RESPONSE TO BROWSER: ", response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	// w.Write([]byte(fmt.Sprintf("%v", response)))
	// c.JSON(http.StatusOK, response)

}

type HelloData struct {
	Rooms []mongodb.Room
}

// HelloWorld will return welcome message for home path
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./template/hello.gohtml"))
	roomService := rooms.NewRoomService()
	rooms, err := roomService.GetRooms()
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusInternalServerError, response)
		return
	}
	// helloData := HelloData{
	// 	Rooms: rooms,
	// }
	// err = t.Execute(w, helloData)
	err = t.Execute(w, rooms)
	if err != nil {
		// panic(err)
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusInternalServerError, response)
		return
	}

	// response := common.ResponseFormatter(http.StatusOK, "success", "load page template successfull", "Hello from MySlack Happy App. A Persistence Chat App... The server is running at the background..")
	// log.Println("RESPONSE TO BROWSER: ", response)
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(response)
	// w.Write([]byte(fmt.Sprintf("%v", response)))
	// c.JSON(http.StatusOK, response)
}
