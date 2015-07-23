package articles

import (
	"github.com/bulletind/khabar-admin/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/bulletind/khabar-admin/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/bulletind/khabar-admin/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/bulletind/khabar-admin/models"
	"gopkg.in/mgo.v2"
)

/**
 * Show
 */

func Show(params martini.Params, r render.Render, db *mgo.Database) {

	article := models.Article{}
	oId := bson.ObjectIdHex(params["_id"])

	err := db.C("articles").FindId(oId).One(&article)

	if err != nil {
		r.HTML(400, "400", err)
	}

	r.HTML(200, "articles/show", article)
}
