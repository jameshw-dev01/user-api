package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jameshw-dev01/user-api/spec"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	Name  string
	Email string
	Age   uint
}

func createUser(c *gin.Context, s ServerContext) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	var userResponse UserResponse
	c.BindJSON(&userResponse)
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
	}

}

func updateUser(c *gin.Context, s ServerContext) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	bcrypt.CompareHashAndPassword([]byte(password), []byte(s.Users[username].Hash))
}

func getUser(c *gin.Context, s ServerContext) {

}
