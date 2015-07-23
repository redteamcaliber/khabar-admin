package available

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/bulletind/khabar-admin/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

/**
 * Get add/edit form
 */

func AddEdit(r render.Render) {
	r.HTML(200, "available/form", nil)
}

/**
 * Add
 */

func Add(topic models.AvailableTopic, r render.Render, db *mgo.Database) {

	topic.BeforeSave()
	err := db.C("topics_available").Insert(topic)

	if err != nil {
		r.HTML(400, "400", err)
	} else {
		r.Redirect("/available")
	}
}

/**
 * List
 */

func List(r render.Render, params martini.Params, db *mgo.Database) {

	var available []models.AvailableTopic

	err := db.C("topics_available").Find(nil).Sort("-created_on").All(&available)

	if err != nil {
		r.Error(400)
	}

	r.HTML(200, "available/list", available)
}

/**
 * Show
 */

func Show(params martini.Params, r render.Render, db *mgo.Database) {

	topic := models.AvailableTopic{}
	oId := bson.ObjectIdHex(params["_id"])

	err := db.C("topics_available").FindId(oId).One(&topic)

	if err != nil {
		r.HTML(400, "400", err)
	}

	r.HTML(200, "available/show", topic)
}
