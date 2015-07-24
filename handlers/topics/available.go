package available

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/bulletind/khabar-admin/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

const (
	availableCollection = "topics_available"
)

type AvailableForm struct {
	Op       string
	Topic    models.AvailableTopic
	Channels []Channel
}

type Channel struct {
	Name    string
	Enabled bool
}

/**
 * New
 */

func New(r render.Render) {
	available := new(AvailableForm)
	available.Op = "New"
	available.Channels = append(available.Channels,
		Channel{"email", true},
		Channel{"web", true},
		Channel{"push", true},
	)

	r.HTML(200, "available/form", available)
}

/**
 * Create
 */

func Create(topic models.AvailableTopic, r render.Render, db *mgo.Database) {
	topic.BeforeSave()
	err := db.C(availableCollection).Insert(topic)
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
	available := new(AvailableForm)
	available.Op = "Edit"
	available.Channels = append(available.Channels,
		Channel{"email", false},
		Channel{"web", false},
		Channel{"push", false},
	)
	oId := bson.ObjectIdHex(params["_id"])
	err := db.C(availableCollection).FindId(oId).One(&available.Topic)
	if err != nil {
		r.HTML(400, "400", err)
	}

	fmt.Println(available.Channels)

	for i, v := range available.Channels {
		for _, value := range available.Topic.Channels {
			if v.Name == value {
				available.Channels[i].Enabled = true
			}
		}
	}

	r.HTML(200, "available/form", available)
}

/**
 * List
 */

func List(r render.Render, params martini.Params, db *mgo.Database) {
	var available []models.AvailableTopic
	err := db.C(availableCollection).Find(nil).Sort("-updated_on").All(&available)
	if err != nil {
		r.Error(400)
	}
	r.HTML(200, "available/list", available)
}

/**
 * Update
 */

func Update(params martini.Params, topic models.AvailableTopic, r render.Render, db *mgo.Database) {
	query := bson.M{"_id": bson.ObjectIdHex(params["_id"])}
	doc := bson.M{
		"ident":      topic.Ident,
		"app_name":   topic.AppName,
		"channels":   topic.Channels,
		"updated_on": time.Now().UnixNano() / int64(time.Millisecond),
	}
	err := db.C(availableCollection).Update(query, doc)
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
	err := db.C(availableCollection).Remove(query)
	if err != nil {
		r.HTML(400, "400", err)
	} else {
		r.Redirect("/available")
	}
}
