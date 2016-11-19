package mgofun

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//per Mongo convention, do not use camel case
// BaseModel to be emmbered to other struct as audit trail perpurse
type BaseModel struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
	CreatedBy string        `bson:"created_by,omitempty"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty"`
	UpdatedBy string        `bson:"updated_by,omitempty"`
	IsRemoved bool          `bson:"is_removed,omitempty"`
	RemovedAt time.Time     `bson:"removed_at,omitempty"`
	RemovedBy string        `bson:"removed_by,omitempty"`
}
