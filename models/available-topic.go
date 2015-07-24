package models

type AvailableTopic struct {
	BaseModel `bson:",inline"`
	Ident     string   `json:"ident" form:"ident" binding:"required" bson:"ident" required:"true"`
	AppName   string   `json:"app_name" form:"app_name" binding:"required" bson:"app_name" required:"true"`
	Channels  []string `json:"channels" form:"channels" binding:"required" bson:"channels" required:"true"`
}
