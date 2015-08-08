package events

import (
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	kdb "github.com/bulletind/khabar/db"
	"github.com/gin-gonic/gin"
)

// FormEvent holds the event form
type FormEvent struct {
	Event    kdb.AvailableTopic
	Channels []kdb.Channel
}

// New event
func New(c *gin.Context) {
	form := FormEvent{
		Channels: []kdb.Channel{
			kdb.Channel{Name: "email", Enabled: true},
			kdb.Channel{Name: "web", Enabled: true},
			kdb.Channel{Name: "push", Enabled: true},
		},
	}

	c.HTML(http.StatusOK, "events/form", gin.H{
		"title": "New event",
		"form":  form,
	})
}

// Create an event
func Create(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	event := kdb.AvailableTopic{}
	err := c.Bind(&event)
	if err != nil {
		c.Error(err)
		return
	}

	event.PrepareSave()
	event.Ident = sanitize(event.Ident)
	err = db.C(kdb.AvailableTopicCollection).Insert(event)
	if err != nil {
		c.Error(err)
	}
	c.Redirect(http.StatusMovedPermanently, "/events")
}

// Edit an event
func Edit(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	form := FormEvent{
		Channels: []kdb.Channel{
			kdb.Channel{Name: "email", Enabled: false},
			kdb.Channel{Name: "web", Enabled: false},
			kdb.Channel{Name: "push", Enabled: false},
		},
	}
	oID := bson.ObjectIdHex(c.Param("_id"))
	err := db.C(kdb.AvailableTopicCollection).FindId(oID).One(&form.Event)
	if err != nil {
		c.Error(err)
	}

	for i, v := range form.Channels {
		for _, value := range form.Event.Channels {
			if v.Name == value {
				form.Channels[i].Enabled = true
			}
		}
	}

	c.HTML(http.StatusOK, "events/form", gin.H{
		"title": "Edit event",
		"form":  form,
	})
}

// List all events
func List(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	events := []kdb.AvailableTopic{}
	err := db.C(kdb.AvailableTopicCollection).Find(nil).Sort("app_name").All(&events)
	if err != nil {
		c.Error(err)
	}
	c.HTML(http.StatusOK, "events/list", gin.H{
		"title":  "Events",
		"events": events,
	})
}

// Update an event
func Update(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	event := kdb.AvailableTopic{}
	err := c.Bind(&event)
	if err != nil {
		c.Error(err)
		return
	}

	query := bson.M{"_id": bson.ObjectIdHex(c.Param("_id"))}
	event.Ident = sanitize(event.Ident)
	doc := bson.M{
		"ident":      event.Ident,
		"app_name":   event.AppName,
		"channels":   event.Channels,
		"updated_on": time.Now().UnixNano() / int64(time.Millisecond),
	}
	err = db.C(kdb.AvailableTopicCollection).Update(query, doc)
	if err != nil {
		c.Error(err)
	}
	c.Redirect(http.StatusMovedPermanently, "/events")
}

// Delete an event
func Delete(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	query := bson.M{"_id": bson.ObjectIdHex(c.Param("_id"))}
	err := db.C(kdb.AvailableTopicCollection).Remove(query)
	if err != nil {
		c.Error(err)
	}
	c.Redirect(http.StatusMovedPermanently, "/events")
}

// sanitize replaces all spaces with underscores
func sanitize(ident string) string {
	// TODO: use regex to replace all other characters other than [A-Za-z0-9_]
	return strings.Replace(strings.Trim(ident, " "), " ", "_", -1)
}
