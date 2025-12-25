package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()

	// 先注册
	//TestRegister(t)

	token := GetToken(router)

	body := `{"title":"Hello","content":"World"}`
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestGetPosts(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
