package available

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	khabar "github.com/bulletind/khabar/db"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type AvailableForm struct {
	Op       string
	Topic    khabar.AvailableTopic
	Channels []khabar.Channel
}

/**
 * New
 */

func New(r render.Render) {
	form := new(AvailableForm)
	form.Op = "New"
	form.Channels = append(
		form.Channels,
		khabar.Channel{Name: "email", Enabled: true},
		khabar.Channel{Name: "web", Enabled: true},
		khabar.Channel{Name: "push", Enabled: true},
	)

	r.HTML(200, "available/form", form)
}

/**
 * Create
 */

func Create(topic khabar.AvailableTopic, r render.Render, db *mgo.Database) {
	topic.PrepareSave()
	topic.Ident = sanitize(topic.Ident)
	err := db.C(khabar.AvailableTopicCollection).Insert(topic)
	if err != nil {
		r.HTML(400, "400", err)
	} else {
		r.Redirect("/available")
	}
}

/**
 * Edit
 */

func Edit(params martini.Params, r render.Render, db *mgo.Database) {
	form := new(AvailableForm)
	form.Op = "Edit"
	form.Channels = append(form.Channels,
		khabar.Channel{Name: "email", Enabled: false},
		khabar.Channel{Name: "web", Enabled: false},
		khabar.Channel{Name: "push", Enabled: false},
	)
	oId := bson.ObjectIdHex(params["_id"])
	err := db.C(khabar.AvailableTopicCollection).FindId(oId).One(&form.Topic)
	if err != nil {
		r.HTML(400, "400", err)
	}

	for i, v := range form.Channels {
		for _, value := range form.Topic.Channels {
			if v.Name == value {
				form.Channels[i].Enabled = true
			}
		}
	}

	r.HTML(200, "available/form", form)
}

/**
 * List
 */

func List(r render.Render, params martini.Params, db *mgo.Database) {
	var available []khabar.AvailableTopic
	err := db.C(khabar.AvailableTopicCollection).Find(nil).Sort("-updated_on").All(&available)
	if err != nil {
		r.Error(400)
	}
	r.HTML(200, "available/list", available)
}

/**
 * Update
 */

func Update(params martini.Params, topic khabar.AvailableTopic, r render.Render, db *mgo.Database) {
	query := bson.M{"_id": bson.ObjectIdHex(params["_id"])}

	topic.Ident = sanitize(topic.Ident)
	doc := bson.M{
		"ident":      topic.Ident,
		"app_name":   topic.AppName,
		"channels":   topic.Channels,
		"updated_on": time.Now().UnixNano() / int64(time.Millisecond),
	}
	err = db.C(khabar.AvailableTopicCollection).Update(query, doc)
	if err != nil {
		r.HTML(400, "400", err)
	} else {
		r.Redirect("/available")
	}
}

/**
 * Delete
 */

func Delete(params martini.Params, r render.Render, db *mgo.Database) {
	query := bson.M{"_id": bson.ObjectIdHex(params["_id"])}
	err := db.C(khabar.AvailableTopicCollection).Remove(query)
	if err != nil {
		r.HTML(400, "400", err)
	} else {
		r.Redirect("/available")
	}
}

/**
 * Sanitize
 */

func sanitize(ident string) string {
	// TODO: use regex to replace all other characters other than [A-Za-z0-9_]
	return strings.Replace(ident, " ", "_", -1)
}
