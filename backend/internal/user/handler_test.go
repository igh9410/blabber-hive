package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) IsUserRegistered(c context.Context, email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockService) CreateUser(c context.Context, req *CreateUserReq, email string) (*CreateUserRes, error) {
	args := m.Called(c, req, email)
	return args.Get(0).(*CreateUserRes), args.Error(1)

}

/*
func TestHandleOAuth2Callback(t *testing.T) {
	// Initialize Gin engine
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockService := new(MockService)
	h := NewHandler(mockService)
	r.GET("/oauth2callback", h.HandleOAuth2Callback)

	t.Run("User not registered", func(t *testing.T) {
		mockService.On("IsUserRegistered", testEmail).Return(false, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/oauth2callback", nil)
		req = req.WithContext(context.WithValue(req.Context(), "email", ))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/register", w.Header().Get("Location"))
	})

	// ... other test cases for HandleOAuth2Callback
} */

func TestCreateUserHandler(t *testing.T) {

	mockService := new(MockService)
	h := NewHandler(mockService)
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/api/users/register", h.CreateUser)

	userReq := CreateUserReq{
		// Populate fields
		Username:        testUsername,
		ProfileImageURL: nil,
	}
	mockService.On("CreateUser", mock.Anything, userReq, testEmail).Return(&CreateUserRes{}, nil)
	jsonData, err := json.Marshal(userReq)
	if err != nil {
		t.Fatalf("Failed to marshal user request: %v", err)
	}
	req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json") // set content type to JSON

	c := gin.Context{Request: req}
	c.Set("email", testEmail)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	fmt.Println("Status = ", w.Result().Status)

	assert.Equal(t, http.StatusOK, w.Code)

}
