package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jameshw-dev01/user-api/spec"
)

type ServerContext struct {
	Users map[string]spec.User
	DB    spec.DbInterface
}

func main() {
	s := ServerContext{}
	router := gin.Default()
	router.POST("/api/v1/user/create", func(c *gin.Context) { createUser(c, s) })
	router.GET("/api/v1/user/:username")
	router.PATCH("/api/v1/user/:username/update")
	router.Run()
}
