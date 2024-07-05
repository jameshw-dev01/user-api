package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jameshw-dev01/user-api/database"
	"github.com/jameshw-dev01/user-api/spec"
)

type ServerContext struct {
	Users map[string]spec.User
	DB    spec.DbInterface
}

func setupRouter(resetDB bool) *gin.Engine {
	s := ServerContext{Users: make(map[string]spec.User), DB: database.GetDBConnection(resetDB, "PROD")}
	users, err := s.DB.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range users {
		s.Users[u.Username] = u
	}
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.POST("/api/v1/user", func(c *gin.Context) { createUser(c, &s) })
	router.GET(
		"/api/v1/user/:username",
		func(c *gin.Context) { runAuth(c, &s) },
		func(c *gin.Context) { getUser(c, &s) },
	)
	router.PUT(
		"/api/v1/user/:username",
		func(c *gin.Context) { runAuth(c, &s) },
		func(c *gin.Context) { updateUser(c, &s) },
	)
	router.DELETE(
		"/api/v1/user/:username",
		func(c *gin.Context) { runAuth(c, &s) },
		func(c *gin.Context) { deleteUser(c, &s) },
	)
	return router
}

func main() {
	router := setupRouter(false)
	router.Run()
}
