package models

import "strings"

type AvailableTopic struct {
	BaseModel `bson:",inline"`
	Ident     string   `json:"ident" form:"ident" binding:"required" bson:"ident" required:"true"`
	AppName   string   `json:"app_name" form:"app_name" binding:"required" bson:"app_name" required:"true"`
	Channels  []string `json:"channels" form:"channels" binding:"required" bson:"channels" required:"true"`
}

func (available *AvailableTopic) Sanitize() {
	// TODO: use regex to replace all other characters other than [A-Za-z0-9_]
	available.Ident = strings.Replace(available.Ident, " ", "_", -1)
}
