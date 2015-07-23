package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type BaseModel struct {
	Id         bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedOn  int64         `json:"created_on" bson:"created_on"`
	ModifiedOn int64         `json:"updated_on" bson:"updated_on"`
}

/**
 * Before any model saves
 */

func (model *BaseModel) BeforeSave() {
	now := time.Now().UnixNano() / int64(time.Millisecond)

	if !model.Id.Valid() {
		model.Id = bson.NewObjectId()
	}

	if !(model.CreatedOn > 0) {
		model.CreatedOn = now
	}

	model.ModifiedOn = now
}

func (model *BaseModel) Date() time.Time {
	return time.Unix(model.CreatedOn, 0)
}
