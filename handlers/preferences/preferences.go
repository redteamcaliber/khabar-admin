package preferences

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	kdb "github.com/bulletind/khabar/db"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// FormPreferences is a struct for preferences form
type FormPreferences struct {
	Preferences []kdb.Topic
	Channels    []string
}

// GlobalPreference holds the default preferences for users and orgs
type GlobalPreference struct {
	Preference []kdb.Topic
}

// List responds with list of global notification preferences
func List(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var global []kdb.Topic
	var events []kdb.AvailableTopic

	preferences := []kdb.Topic{}
	form := FormPreferences{
		Channels: []string{"email", "web", "push"},
	}
	query := bson.M{
		"org":  "",
		"user": "",
	}

	err := db.C(kdb.TopicCollection).Find(query).Sort("-updated_on").All(&global)
	if err != nil {
		c.Error(err)
	}

	// Get all events
	err = db.C(kdb.AvailableTopicCollection).Find(nil).Sort("app_name", "sortindex").All(&events)
	if err != nil {
		c.Error(err)
	}

	for _, event := range events {
		channels := []kdb.Channel{}
		for _, channel := range event.Channels {
			channels = append(channels, kdb.Channel{
				Name:    channel,
				Locked:  false,
				Default: false,
			})
		}
		preferences = append(preferences, kdb.Topic{
			Ident:    event.Ident,
			Channels: channels,
		})
	}

	for _, event := range events {
		for _, channel := range event.Channels {

			for _, pref := range global {
				for _, ch := range pref.Channels {

					for i, p := range preferences {
						for j, c := range p.Channels {

							if ch.Name == channel && pref.Ident == event.Ident && c.Name == channel && p.Ident == event.Ident {
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

	c.HTML(http.StatusOK, "preferences/form", gin.H{
		"title": "Global notification preferences",
		"form":  form,
	})
}

// Update updates the global notification preferences
func Update(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	err := c.Request.ParseForm()
	if err != nil {
		c.Error(err)
	}

	global := GlobalPreference{}
	err = decoder.Decode(&global, c.Request.PostForm)
	if err != nil {
		c.Error(err)
	}

	preferences := global.Preference

	// Update
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
		_, err := db.C(kdb.TopicCollection).Upsert(query, doc)

		if err != nil {
			log.Println(err)
			break
		}
	}

	c.Redirect(http.StatusMovedPermanently, "/preferences")
}
