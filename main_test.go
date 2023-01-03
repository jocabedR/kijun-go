package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestExistentUserRoute(t *testing.T) {

	body := gin.H{
		"Id":                14,
		"Username":          "jnovichenkore",
		"Name":              "Jacqueline Novichenko",
		"Birth_date":        "1975-02-05T00:00:00Z",
		"Registration_date": "2017-10-13T00:00:00Z",
	}

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/whoami/jnovichenkore", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var user User
	err := json.Unmarshal(w.Body.Bytes(), &user)

	assert.Nil(t, err)
	assert.Equal(t, body["Id"], user.Id)
	assert.Equal(t, body["Username"], user.Username)
	assert.Equal(t, body["Name"], user.Name)
	assert.Equal(t, body["Birth_date"], user.Birth_date)
	assert.Equal(t, body["Registration_date"], user.Registration_date)
}

func TestNonExistentUserRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/whoami/jnovichenkor", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	var message map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &message)
	value, exists := message["message"]

	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, "User not found", value)
}
