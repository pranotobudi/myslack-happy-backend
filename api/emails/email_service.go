package emails

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/pranotobudi/myslack-happy-backend/api/messages"
	"github.com/pranotobudi/myslack-happy-backend/api/users"
	"github.com/pranotobudi/myslack-happy-backend/common"
	"github.com/pranotobudi/myslack-happy-backend/mongodb"
)

type IEmailService interface {
	MailChat(userMongo mongodb.User) (string, error)
}

type emailService struct {
	repo mongodb.IMongoDB
}
type EmailChat struct {
	Email string
	Rooms []EmailRoom
}

func (e EmailChat) String() string {
	return fmt.Sprintf("Email: %v\n Rooms:\n%v\n", e.Email, e.Rooms)
}

type EmailRoom struct {
	RoomId   string
	Messages []mongodb.Message
}

func (e EmailRoom) String() string {
	return fmt.Sprintf("RoomId: %v\n Messages:\n%v\n", e.RoomId, e.Messages)
}

// NewEmailService will return emailService object
func NewEmailService() *emailService {
	r := mongodb.NewMongoDB()
	return &emailService{repo: r}
}

// UserMailChat will send email of chat with each room and its messages to the current user
func (s *emailService) MailChat(userMongo mongodb.User) (string, error) {
	// remove all user rooms first
	// update := bson.D{{"$set", bson.M{"rooms": []string{}}}}
	// err := s.repo.UpdateUser(filter, update, opts)
	// if err != nil {
	// 	return nil, err
	// }
	// paths := []string{
	// 	"./template/email.html",
	// }
	messageService := messages.NewMessageService()
	userService := users.NewUserService()
	user, err := userService.GetUser(userMongo.Email)
	if err != nil {
		return "", err
	}

	// Load EmailChat data
	emailChat := EmailChat{}
	emailChat.Email = user.Email

	for _, roomId := range user.Rooms {
		message, err := messageService.GetMessages(roomId)
		if err != nil {
			return "", err
		}
		emailRoom := EmailRoom{}
		emailRoom.RoomId = roomId
		emailRoom.Messages = message

		emailChat.Rooms = append(emailChat.Rooms, emailRoom)
	}
	// log.Println("EMAILCHAT STRUCT: ", emailChat)

	// t := template.Must(template.New("html-tmpl").ParseFiles(paths...))
	t := template.Must(template.ParseFiles("./template/email.gohtml"))
	buf := new(bytes.Buffer)
	// err = t.Execute(os.Stdout, emailChat)
	err = t.Execute(buf, emailChat)
	if err != nil {
		// panic(err)
		log.Println("LOADING TEMPLATE ERROR.")
	}
	toEmail := []string{emailChat.Email}
	common.SendEmail(toEmail, buf.String())

	return "email sent succesfully to your email", nil
}
