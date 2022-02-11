package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pranotobudi/myslack-happy-backend/api/messages"
	"github.com/pranotobudi/myslack-happy-backend/api/rooms"
	"github.com/pranotobudi/myslack-happy-backend/api/users"
	"github.com/pranotobudi/myslack-happy-backend/config"
	"github.com/pranotobudi/myslack-happy-backend/msgserver"
)

func main() {
	StartApp()
}
func StartApp() {
	if os.Getenv("APP_ENV") != "production" {
		// executed in development, because we need to load .env variables in local env
		// in development only,
		// load local env variables to os
		// for production set those OS environment on production environment settings
		// production env like heroku provide that "APP_ENV" variable
		// for other platform (kubernetes): set APP_ENV = production manually
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("failed to load .env file")
		}
		log.Println("Load development environment variables..")
	}

	// # run router server
	gin.SetMode(gin.ReleaseMode)
	appConfig := config.AppConfig()
	log.Println("server run on port:8080...")
	http.ListenAndServe(":"+appConfig.Port, Router())

}
func Router() *gin.Engine {
	// handler
	messageHandler := messages.NewMessageHandler()
	roomHandler := rooms.NewRoomHandler()
	userHandler := users.NewUserHandler()

	// #2 init gin main server
	router := gin.Default()
	router.Use(CORS)
	// gin.SetMode(gin.ReleaseMode)

	// #3 handle url to init websocket client connection (will have func to handle incoming url)
	// this client will notify subscribe event to the global message server through channel.

	router.GET("/", users.HelloWorld)
	router.GET("/rooms", roomHandler.GetRooms)
	router.POST("/room", roomHandler.AddRoom)
	router.GET("/room", roomHandler.GetAnyRoom)
	router.GET("/messages", messageHandler.GetMessages)
	router.GET("/userByEmail", userHandler.GetUserByEmail)
	router.POST("/userAuth", userHandler.UserAuth)
	router.PUT("/updateUserRooms", userHandler.UpdateUserRooms)
	router.GET("/websocket", msgserver.InitWebsocket)

	return router
}

func CORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

// func serveMainPage(c *gin.Context) {
// 	// request: userId
// 	// response: user snapshot to load main page
// 	fmt.Println("inside serveMainPage!")
// 	c.File("static/index.html")
// }
// func serveStaticPage(c *gin.Context) {
// 	fmt.Println("inside serveStaticPage!")
// 	filePath := c.Request.URL.Path
// 	c.File(filePath)
// }

// func chatServer(c *gin.Context) {
// 	log.Println("inside chatServer! message Send..")
// 	c.Writer.WriteHeader(http.StatusAccepted)
// 	c.Writer.Write([]byte("msg send.."))
// }
