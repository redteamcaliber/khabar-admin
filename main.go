package main

import (
	"github.com/bulletind/khabar-admin/handlers/topics"
	"github.com/bulletind/khabar-admin/middlewares"
	"github.com/bulletind/khabar-admin/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

/**
 * Main
 *
 * - use middlewares
 * - define routes
 * - listen and serve
 */

func main() {

	// Initialize
	m := martini.Classic()

	// Connect to mongo
	m.Use(middlewares.Connect())

	// Templating support
	m.Use(middlewares.Templates())

	// Routes
	m.Group("/topics", func(r martini.Router) {
		r.Get("", topics.List)
		r.Get("/new", topics.AddEdit)
		r.Get("/:_id", topics.Show)
		r.Post("", binding.Bind(models.AvailableTopic{}), topics.Add)
	})

	// Start listening
	m.Run()
}
