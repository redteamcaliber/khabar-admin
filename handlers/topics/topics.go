package topics

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	khabar "github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/available_topics"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type TopicForm struct {
	Op       string
	Topic    khabar.Topic
	Channels []khabar.Channel
}

/**
 * List
 */

func List(r render.Render, params martini.Params, db *mgo.Database) {
	var global []khabar.Topic
	var available []khabar.AvailableTopic
	preferences := map[string]available_topics.ChotaTopic{}
	query := bson.M{
		"org":  "",
		"user": "",
	}
	err := db.C(khabar.TopicCollection).Find(query).Sort("-updated_on").All(&global)
	if err != nil {
		r.Error(400)
	}

	// Get the available topics

	err = db.C(khabar.AvailableTopicCollection).Find(nil).Sort("-updated_on").All(&available)
	if err != nil {
		r.Error(400)
	}

	// Remove the non-available ones

	for _, availableTopic := range available {
		ct := available_topics.ChotaTopic{}
		for _, channel := range availableTopic.Channels {
			ct[channel] = &available_topics.TopicDetail{Locked: false, Default: false}
		}
		preferences[availableTopic.Ident] = ct
	}

	for _, topic := range available {
		for _, channel := range topic.Channels {
			// Remove it from the global
			for _, pref := range global {
				for _, ch := range pref.Channels {
					if ch.Name == channel && pref.Ident == topic.Ident {
						preferences[topic.Ident][channel].Default = ch.Default
						preferences[topic.Ident][channel].Locked = ch.Locked
					}
				}
			}
		}
	}

	r.HTML(200, "topics/form", preferences)
}

/**
 * Update
 */

func Update(params martini.Params, req *http.Request, r render.Render, db *mgo.Database) {

	req.ParseForm()

	r.Redirect("/topics")
}
