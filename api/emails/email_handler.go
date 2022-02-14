package emails

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pranotobudi/myslack-happy-backend/common"
	"github.com/pranotobudi/myslack-happy-backend/mongodb"
)

type IEmailHandler interface {
	MailChat(w http.ResponseWriter, r *http.Request)
}
type emailHandler struct {
	emailService IEmailService
}

// NewEmailHandler will initialize emailHandler object
func NewEmailHandler() *emailHandler {
	emailService := NewEmailService()
	return &emailHandler{emailService: emailService}
}

func (h *emailHandler) MailChat(w http.ResponseWriter, r *http.Request) {
	// login
	var userMongo mongodb.User
	err := json.NewDecoder(r.Body).Decode(&userMongo)
	// err := c.BindJSON(&userMongo)
	log.Println("UserMailChat userMongo: ", userMongo)

	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusBadRequest, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusBadRequest, err)
		return
	}
	msg, err := h.emailService.MailChat(userMongo)
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		// w.Write([]byte(fmt.Sprintf("%v", response)))
		// c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := common.ResponseFormatter(http.StatusOK, "success", "email sent successfully", msg)
	log.Println("RESPONSE TO BROWSER: ", response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	// w.Write([]byte(fmt.Sprintf("%v", response)))
	// c.JSON(http.StatusOK, response)
}
