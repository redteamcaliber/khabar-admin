package main

import (
	"github.com/bulletind/khabar-admin/handlers/topics"
	"github.com/bulletind/khabar-admin/middlewares"
	"github.com/bulletind/khabar-admin/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
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

	m.Get("/", func(r render.Render) {
		r.Redirect("/available")
	})

	m.Group("/available", func(r martini.Router) {
		r.Get("", available.List)
		r.Get("/new", available.AddEdit)
		r.Get("/:_id", available.Show)
		r.Post("", binding.Bind(models.AvailableTopic{}), available.Add)
	})

	// Start listening
	m.Run()
}
