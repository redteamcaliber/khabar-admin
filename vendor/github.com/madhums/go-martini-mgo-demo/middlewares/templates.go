package middlewares

import (
	"github.com/bulletind/khabar-admin/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/bulletind/khabar-admin/Godeps/_workspace/src/github.com/martini-contrib/render"
)

/**
 * Set templating defaults
 */

func Templates() martini.Handler {
	return render.Renderer(render.Options{
		Layout: "layouts/default",
	})
}
