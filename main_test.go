package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	main "github.com/pranotobudi/myslack-happy-backend"
	"github.com/stretchr/testify/assert"
)

// TestRouting will test each end point is active
func TestRouting(t *testing.T) {
	tt := []struct {
		Name   string
		Path   string
		Method string
		Status int
	}{
		{"Get home", "/", http.MethodGet, http.StatusOK},
		{"Get rooms", "/rooms", http.MethodGet, http.StatusOK},
		// {"Get room", "/room", http.MethodGet, http.StatusOK},
		// {"Post room", "/room", http.MethodPost, http.StatusOK},
		// {"Get messages", "/messages", http.MethodGet, http.StatusOK},
		// {"Get userByemail", "/userByEmail", http.MethodGet, http.StatusOK},
		// {"Post userAuth", "/userAuth", http.MethodPost, http.StatusOK},
		// {"Put updateUserRooms", "/updateUserRooms", http.MethodPut, http.StatusOK},
		// {"Get websocket", "/websocket", http.MethodGet, http.StatusOK},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest(tc.Method, tc.Path, nil)
			w := httptest.NewRecorder()

			router := main.Router()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
