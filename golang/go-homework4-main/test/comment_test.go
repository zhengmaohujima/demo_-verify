package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()
	token := GetToken(router)

	body := `{"content":"Nice post!"}`
	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}
