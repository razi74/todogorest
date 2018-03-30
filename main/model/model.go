package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Task struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Taskname string        `bson:"taskname" json:"taskname"`
}
