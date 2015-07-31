package main

import (
	"github.com/bulletind/khabar-admin/handlers/available"
	"github.com/bulletind/khabar-admin/handlers/topics"
	"github.com/bulletind/khabar-admin/middlewares"
	khabar "github.com/bulletind/khabar/db"
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
		r.Get("/new", available.New)
		r.Get("/:_id", available.Edit)
		r.Post("", binding.Bind(khabar.AvailableTopic{}), available.Create)
		r.Post("/:_id", binding.Bind(khabar.AvailableTopic{}), available.Update)
		r.Post("/:_id/delete", available.Delete)
	})

	m.Group("/topics", func(r martini.Router) {
		r.Get("", topics.List)
		r.Post("", topics.Update)
	})

	// Start listening
	m.Run()
}
