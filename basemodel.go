package mgofun

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// seperate by _
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

// Camel case
// BaseModel to be emmbered to other struct as audit trail perpurse
type BaseModelCamel struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	CreatedAt time.Time     `bson:"createdAt,omitempty"`
	CreatedBy string        `bson:"createdBy,omitempty"`
	UpdatedAt time.Time     `bson:"updatedAt,omitempty"`
	UpdatedBy string        `bson:"updatedBy,omitempty"`
	IsRemoved bool          `bson:"isRemoved,omitempty"`
	RemovedAt time.Time     `bson:"removedAt,omitempty"`
	RemovedBy string        `bson:"removedBy,omitempty"`
}
