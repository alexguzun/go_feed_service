package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type FeedEntry struct {
	ID       string        `bson:"_id,omitempty"`
	Title    string        `bson:"title,omitempty"`
	SourceID bson.ObjectId `bson:"source_id,omitempty"`
	Content  string        `bson:"content,omitempty"`
	AddedOn  time.Time     `bson:"added_on,omitempty"`
}
