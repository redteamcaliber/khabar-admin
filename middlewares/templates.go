package middlewares

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

/**
 * Set templating defaults
 */

func Templates() martini.Handler {
	return render.Renderer(render.Options{
		Layout: "layouts/default",
		// Funcs: []template.FuncMap{
		// 	{
		// 		"formatTime": func(args ...interface{}) string {
		// 			t1 := time.Unix(args[0].(int64), 0)
		// 			return t1.Format(time.Stamp)
		// 		},
		// 	},
		// },
	})
}
