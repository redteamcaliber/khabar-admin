package topics

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	khabar "github.com/bulletind/khabar/db"
	"github.com/go-martini/martini"
	"github.com/gorilla/schema"
	"github.com/martini-contrib/render"
)

var decoder = schema.NewDecoder()

type PreferencesForm struct {
	Preferences []khabar.Topic
	Channels    []string
}

type GlobalPreference struct {
	Preference []khabar.Topic
}

/**
 * List
 */

func List(r render.Render, params martini.Params, db *mgo.Database) {
	var form = PreferencesForm{
		Channels: []string{"email", "web", "push"},
	}
	var global []khabar.Topic
	var available []khabar.AvailableTopic
	preferences := []khabar.Topic{}

	query := bson.M{
		"org":  "",
		"user": "",
	}
	err := db.C(khabar.TopicCollection).Find(query).Sort("-updated_on").All(&global)
	if err != nil {
		r.Error(400)
	}

	// Get the available topics

	err = db.C(khabar.AvailableTopicCollection).Find(nil).Sort("app_name").All(&available)
	if err != nil {
		r.Error(400)
	}

	for _, availableTopic := range available {
		channels := []khabar.Channel{}
		for _, channel := range availableTopic.Channels {
			channels = append(channels, khabar.Channel{
				Name:    channel,
				Locked:  false,
				Default: false,
			})
		}
		preferences = append(preferences, khabar.Topic{
			Ident:    availableTopic.Ident,
			Channels: channels,
		})
	}

	for _, topic := range available {
		for _, channel := range topic.Channels {

			for _, pref := range global {
				for _, ch := range pref.Channels {

					for i, p := range preferences {
						for j, c := range p.Channels {

							if ch.Name == channel && pref.Ident == topic.Ident && c.Name == channel && p.Ident == topic.Ident {
								preferences[i].Channels[j].Default = ch.Default
								preferences[i].Channels[j].Locked = ch.Locked
							}

						}
					}

				}
			}

		}
	}

	form.Preferences = preferences

	r.HTML(200, "topics/form", form)
}

/**
 * Update
 */

func Update(params martini.Params, req *http.Request, r render.Render, db *mgo.Database) {

	// Parse the form
	err := req.ParseForm()

	if err != nil {
		r.Error(400)
	}

	// Decode it back to GlobalPreference struct
	globalPreferences := new(GlobalPreference)
	err = decoder.Decode(globalPreferences, req.PostForm)

	if err != nil {
		r.Error(400)
	}

	// Update
	preferences := globalPreferences.Preference

	for _, preference := range preferences {
		query := bson.M{
			"org":   "",
			"user":  "",
			"ident": preference.Ident,
		}
		doc := bson.M{
			"$set": bson.M{
				"channels": preference.Channels,
			},
		}
		_, err := db.C(khabar.TopicCollection).Upsert(query, doc)

		if err != nil {
			log.Println(err)
			break
		}
	}

	r.Redirect("/topics")
}
