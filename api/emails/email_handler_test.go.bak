package emails

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pranotobudi/myslack-happy-backend/common"
	"github.com/pranotobudi/myslack-happy-backend/mongodb"
	"github.com/stretchr/testify/assert"
)

var (
	mailChatServiceFunc func(userMongo mongodb.User) (string, error)
)

type mockEmailService struct{}

func (m *mockEmailService) MailChat(userMongo mongodb.User) (string, error) {
	return mailChatServiceFunc(userMongo)
}

func TestGetMessagesHandler(t *testing.T) {

	tt := []struct {
		Name     string
		mockFunc func(userMongo mongodb.User) (string, error)
		CodeWant int
		Body     []byte
	}{
		{
			Name: "MailChat Failed BadRequest",
			mockFunc: func(userMongo mongodb.User) (string, error) {
				return "", errors.New("fail to decode JSON")
			},
			CodeWant: http.StatusBadRequest,
			Body:     []byte(``),
		},
		{
			Name: "MailChat Failed",
			mockFunc: func(userMongo mongodb.User) (string, error) {
				return "", errors.New("fail to get messages")
			},
			CodeWant: http.StatusInternalServerError,
			Body:     []byte(`{"id":"abc123", "email":"budi@gmail.com", "username":"budi", "user_image":"https:aws.com", "rooms":["61cc50877ea033031b1a950e"]}`),
		},
		{
			Name: "MailChat Success",
			mockFunc: func(userMongo mongodb.User) (string, error) {
				return "Email has been sent successfully", nil
			},
			CodeWant: http.StatusOK,
			Body:     []byte(`{"id":"abc123", "email":"budi@gmail.com", "username":"budi", "user_image":"https:aws.com", "rooms":["61cc50877ea033031b1a950e"]}`),
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			mailChatServiceFunc = tc.mockFunc

			emailHandler := NewEmailHandler()
			emailHandler.emailService = &mockEmailService{}
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/mailChat", bytes.NewBuffer(tc.Body))

			log.Println(req.RequestURI)
			emailHandler.MailChat(rr, req)

			log.Println("test response: ", rr.Body.String())
			// check header StatusCode
			assert.EqualValues(t, tc.CodeWant, rr.Code)
			// check response (JSON format) StatusCode
			var response common.Response
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				assert.Errorf(t, err, "response format is not valid")
			}
			assert.EqualValues(t, tc.CodeWant, response.Meta.Code)
		})
	}

}
