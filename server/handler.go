package main

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jameshw-dev01/user-api/spec"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   uint   `json:"age"`
}

func isUserValid(user UserResponse) bool {
	// Matches <START>anystring@anystring.anystring<END>
	emailRegex := regexp.MustCompile(`\S+@\S+\.\S+`)
	return user.Name != "" && emailRegex.MatchString(user.Email)
}

func createUser(c *gin.Context, s *ServerContext) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("missing or malformed basic auth header"))
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var userResponse UserResponse
	err = c.BindJSON(&userResponse)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if !isUserValid(userResponse) {
		c.AbortWithError(http.StatusBadRequest, errors.New("user data is invalid"))
		return
	}
	if _, keyFound := s.Users[username]; keyFound {
		c.AbortWithError(http.StatusBadRequest, errors.New("username already in use"))
		return
	}
	user := spec.User{
		Username: username,
		Hash:     string(hash),
		Email:    userResponse.Email,
		Name:     userResponse.Name,
		Age:      userResponse.Age,
	}
	err = s.DB.Create(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		s.Users[username] = user
		c.IndentedJSON(http.StatusCreated, userResponse)
	}
}

func runAuth(c *gin.Context, s *ServerContext) {
	requestedUsername := c.Param("username")
	username, password, ok := c.Request.BasicAuth()
	if requestedUsername != username {
		c.AbortWithError(http.StatusBadRequest, errors.New("username and auth do not match"))
		return
	}
	if s.Users[username].Username == "" {
		c.AbortWithError(http.StatusNotFound, errors.New("username not found"))
		return
	}
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("failed to parse auth header"))
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(s.Users[username].Hash), []byte(password))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func updateUser(c *gin.Context, s *ServerContext) {
	username := c.Param("username")
	var userResponse UserResponse
	err := c.BindJSON(&userResponse)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user := spec.User{
		Username: username,
		Hash:     s.Users[username].Hash,
		Email:    userResponse.Email,
		Name:     userResponse.Name,
		Age:      userResponse.Age,
	}
	err = s.DB.Update(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		s.Users[username] = user
		c.IndentedJSON(http.StatusOK, userResponse)
	}
}

func getUser(c *gin.Context, s *ServerContext) {
	username := c.Param("username")
	user := s.Users[username]
	if user.Username == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userResponse := UserResponse{Name: user.Name, Email: user.Email, Age: user.Age}
	c.IndentedJSON(http.StatusOK, userResponse)
}

func deleteUser(c *gin.Context, s *ServerContext) {
	username := c.Param("username")
	user := spec.User{
		Username: username,
		Hash:     s.Users[username].Hash,
		Email:    s.Users[username].Email,
		Name:     s.Users[username].Name,
		Age:      s.Users[username].Age,
	}
	err := s.DB.Delete(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		delete(s.Users, username)
		c.IndentedJSON(http.StatusOK, s.Users[username])
	}
}
