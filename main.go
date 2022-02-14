package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/pranotobudi/myslack-happy-backend/api/emails"
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
	// gin.SetMode(gin.ReleaseMode)
	appConfig := config.AppConfig()
	log.Println("server run on port:8080...")
	http.ListenAndServe(":"+appConfig.Port, Router())

}
func Router() *chi.Mux {
	// handler
	messageHandler := messages.NewMessageHandler()
	roomHandler := rooms.NewRoomHandler()
	userHandler := users.NewUserHandler()
	emailHandler := emails.NewEmailHandler()
	// #2 init chi routing server
	router := chi.NewRouter()
	// router := gin.Default()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// #3 handle url to init websocket client connection (will have func to handle incoming url)
	// this client will notify subscribe event to the global message server through channel.
	router.Get("/", users.HelloWorld)
	router.Get("/rooms", roomHandler.GetRooms)
	router.Post("/room", roomHandler.AddRoom)
	router.Get("/room", roomHandler.GetAnyRoom)
	router.Get("/messages", messageHandler.GetMessages)
	router.Get("/userByEmail", userHandler.GetUserByEmail)
	router.Post("/userAuth", userHandler.UserAuth)
	router.Post("/mailChat", emailHandler.MailChat)
	router.Put("/updateUserRooms", userHandler.UpdateUserRooms)
	router.Get("/websocket", msgserver.InitWebsocket)

	return router
}
