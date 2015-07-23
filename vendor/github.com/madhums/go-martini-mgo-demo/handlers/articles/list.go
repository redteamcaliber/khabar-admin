package articles

import (
	"github.com/bulletind/khabar-admin/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/bulletind/khabar-admin/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/bulletind/khabar-admin/models"
	"gopkg.in/mgo.v2"
)

/**
 * List
 */

func List(r render.Render, params martini.Params, db *mgo.Database) {

	var articles []models.Article

	err := db.C("articles").Find(nil).All(&articles)

	if err != nil {
		r.Error(400)
	}

	r.HTML(200, "articles/list", articles)
}
