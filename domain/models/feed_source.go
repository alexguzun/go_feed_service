package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type FeedSource struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	RssUrls     []string      `bson:"rss_urls,omitempty"`
	Description string        `bson:"description,omitempty"`
	Active      bool          `bson:"active,omitempty"`
	UseFeedGuid bool          `bson:"use_feed_guid,omitempty"`
	LastCheck   time.Time     `bson:"last_check,omitempty"`
}
