package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

// FindUserByEmail implements Service.
func (*MockService) FindUserByEmail(ctx context.Context, email string) (*UserDTO, error) {
	panic("unimplemented")
}

func (m *MockService) IsUserRegistered(c context.Context, email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockService) CreateUser(c context.Context, req *CreateUserReq, userID uuid.UUID, email string) (*CreateUserRes, error) {
	args := m.Called(c, req, userID, email)
	return args.Get(0).(*CreateUserRes), args.Error(1)

}

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
	mockService.On("CreateUser", mock.Anything, userReq, testUserID, testEmail).Return(&CreateUserRes{}, nil)

	jsonData, err := json.Marshal(userReq)
	if err != nil {
		t.Fatalf("Failed to marshal user request: %v", err)
	}
	req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json") // set content type to JSON

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user_id", testUserID)
	c.Set("email", testEmail)
	log.Println("testUserID = ", testUserID)
	log.Println("testEmail = ", testEmail)
	log.Println("Context ID = ", c.Value("user_id"))
	log.Println("Context Email = ", c.Value("email"))

	r.ServeHTTP(w, req)
	fmt.Println("Body = ", w.Body)
	fmt.Println("Status = ", w.Result().StatusCode)

	//	assert.Equal(t, http.StatusOK, w.Code)

}
