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
type BaseMod struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	CreatedAt time.Time     `bson:"CreatedAt,omitempty"`
	CreatedBy string        `bson:"CreatedBy,omitempty"`
	UpdatedAt time.Time     `bson:"UpdatedAt,omitempty"`
	UpdatedBy string        `bson:"UpdatedBy,omitempty"`
	IsRemoved bool          `bson:"IsRemoved,omitempty"`
	RemovedAt time.Time     `bson:"RemovedAt,omitempty"`
	RemovedBy string        `bson:"RemovedBy,omitempty"`
}

//ChangeLog
type ChangeLog struct {
	BaseMod      `bson:",inline"`
	ModelObjId   bson.ObjectId `bson:"ModelObjId,omitempty"`
	ModelName    string        `bson:"ModelName,omitempty"`
	ModelValue   interface{}   `bson:"ModelValue,omitempty"`
	ChangeReason string        `bson:"ChangeReason,omitempty"`
}
