// Package main is the CLI.
// You can use the CLI via Terminal.
package main

import (
	"net/http"
	"os"

	"github.com/bulletind/khabar-admin/db"
	"github.com/bulletind/khabar-admin/handlers/events"
	"github.com/bulletind/khabar-admin/handlers/preferences"
	"github.com/bulletind/khabar-admin/middlewares"
	"github.com/gin-gonic/gin"
)

const (
	// Port at which the server starts listening
	Port = "7000"
)

func init() {
	db.Connect()
}

func main() {

	// Configure
	router := gin.Default()
	router.HTMLRender = render()
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true

	// Middlewares
	router.Use(middlewares.Connect)
	router.Use(middlewares.ErrorHandler)

	// Statics
	router.Static("/public", "./public")

	// Routes

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/events")
	})

	// Events
	router.GET("/new/events", events.New)
	router.GET("/events/:_id", events.Edit)
	router.GET("/events", events.List)
	router.POST("/events", events.Create)
	router.POST("/events/:_id", events.Update)
	router.POST("/delete/events/:_id", events.Delete)

	// Preferences
	router.GET("/preferences", preferences.List)
	router.POST("/preferences", preferences.Update)

	// Start listening
	port := Port
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	router.Run(":" + port)
}
