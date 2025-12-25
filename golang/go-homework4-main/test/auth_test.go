package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()

	body := `{
		"username": "testuser",
		"password": "123456",
		"email": "test@test.com"
	}`

	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}
func TestLogin(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()

	body := `{
		"username": "testuser",
		"password": "123456"
	}`

	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
