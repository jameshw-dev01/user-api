package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailValid(t *testing.T) {
	assert.True(t, isUserValid(UserResponse{Name: "John Doe", Email: "test@example.com", Age: 24}))
}

func TestPostUser(t *testing.T) {
	router := setupRouter(true)
	w := httptest.NewRecorder()
	user := UserResponse{Name: "John Doe", Email: "test@example.com", Age: 24}
	jsonUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(jsonUser)))
	req.SetBasicAuth("john_doe", "pass123")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var retrievedUser UserResponse
	json.Unmarshal(w.Body.Bytes(), &retrievedUser)
	assert.Equal(t, user, retrievedUser)
}

func TestGetSuccess(t *testing.T) {
	router := setupRouter(true)
	w := httptest.NewRecorder()
	user := UserResponse{Name: "John Doe", Email: "test@example.com", Age: 24}
	jsonUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(jsonUser)))
	req.SetBasicAuth("john_doe", "pass123")
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/user/john_doe", nil)
	req.SetBasicAuth("john_doe", "pass123")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedUser UserResponse
	json.Unmarshal(w.Body.Bytes(), &retrievedUser)
	assert.Equal(t, user, retrievedUser)
}

func TestGetAuthFail(t *testing.T) {
	jsonUser, _ := json.Marshal(UserResponse{Name: "John Doe", Email: "test@example.com", Age: 24})
	router := setupRouter(true)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(jsonUser)))
	req.SetBasicAuth("john_doe", "pass123")
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/user/john_doe", nil)
	req.SetBasicAuth("john_doe", "PASSFAIL")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetAuthDifferentUser(t *testing.T) {
	user1 := "john_doe"
	pass1 := "pass1"
	user2 := "jane_doe"
	pass2 := "pass2"
	jsonUser1, _ := json.Marshal(UserResponse{Name: "John Doe", Email: "test@example.com", Age: 24})
	jsonUser2, _ := json.Marshal(UserResponse{Name: "Jane Doe", Email: "test1@example.com", Age: 26})
	router := setupRouter(true)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(jsonUser1)))
	req.SetBasicAuth(user1, pass1)
	router.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(jsonUser2)))
	req.SetBasicAuth(user2, pass2)
	router.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/user/"+user1, nil)
	req.SetBasicAuth(user2, pass2)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetNoUser(t *testing.T) {
	user1 := "john_doe"
	user2 := "jane_doe"
	router := setupRouter(true)
	w := httptest.NewRecorder()
	jsonUser, _ := json.Marshal(UserResponse{Name: "John Doe", Email: "test@example.com", Age: 24})

	req, _ := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(jsonUser)))
	req.SetBasicAuth(user1, "pass123")
	router.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/user/"+user2, nil)
	req.SetBasicAuth(user1, "pass123")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSuccess(t *testing.T) {
	user1 := "john_doe"
	userData := UserResponse{Name: "John Doe", Email: "test@example.com", Age: 24}
	pass1 := "pass1"
	router := setupRouter(true)
	w := httptest.NewRecorder()
	jsonUser, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(jsonUser)))
	req.SetBasicAuth(user1, pass1)
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	userData.Age = 25
	userData.Name = "John H Smith"
	userData.Email = "john@example.com"
	jsonUser, _ = json.Marshal(userData)
	// send PUT request
	req, _ = http.NewRequest("PUT", "/api/v1/user/"+user1, strings.NewReader(string(jsonUser)))
	req.SetBasicAuth(user1, pass1)
	router.ServeHTTP(w, req)
	// test PUT responses
	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedUser UserResponse
	json.Unmarshal(w.Body.Bytes(), &retrievedUser)
	assert.Equal(t, userData, retrievedUser)
	// test that it is really saved
	req, _ = http.NewRequest("GET", "/api/v1/user/"+user1, nil)
	req.SetBasicAuth(user1, pass1)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &retrievedUser)
	assert.Equal(t, userData, retrievedUser)
}

func TestDelete(t *testing.T) {
	user1 := "john_doe"
	userData := UserResponse{Name: "John Doe", Email: "test@example.com", Age: 24}
	pass1 := "pass1"
	router := setupRouter(true)
	w := httptest.NewRecorder()
	jsonUser, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(jsonUser)))
	req.SetBasicAuth(user1, pass1)
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	// send DELETE request
	req, _ = http.NewRequest("DELETE", "/api/v1/user/"+user1, strings.NewReader(string(jsonUser)))
	req.SetBasicAuth(user1, pass1)
	router.ServeHTTP(w, req)
	// test DELETE responses
	assert.Equal(t, http.StatusOK, w.Code)
	// test that it is really deleted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/user/"+user1, nil)
	req.SetBasicAuth(user1, pass1)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
