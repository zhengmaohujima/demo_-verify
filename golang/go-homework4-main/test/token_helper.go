package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func GetToken(router http.Handler) string {
	loginBody := `{"username":"testuser","password":"123456"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(loginBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var res map[string]string
	json.Unmarshal(w.Body.Bytes(), &res)
	return res["token"]
}
