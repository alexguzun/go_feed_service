package repository

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alexguzun/go_feed_service/domain/models"
	"github.com/alexguzun/go_feed_service/infrastructure/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type FeedEntriesMongoRepo struct {
	Session *mgo.Session
}

const FeedEntriesCollection string = "feed_entries"

func (repo *FeedEntriesMongoRepo) IsNewEntry(feedEntryId string) (bool, error) {
	count, err := db.Query(repo.Session, func(session *mgo.Session) (interface{}, error) {
		return getEntriesCollection(session).
			Find(bson.M{"_id": feedEntryId}).
			Count()
	})

	if err != nil {
		return false, err
	} else if count.(int) == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (repo *FeedEntriesMongoRepo) Save(newFeedEntry models.FeedEntry) {
	err := db.Execute(repo.Session, func(session *mgo.Session) (err error) {
		return getEntriesCollection(session).Insert(newFeedEntry)
	})

	if err != nil {
		log.WithField("feed_entry_id", newFeedEntry.ID).WithError(err).Error("Failed to add new feed entry")
	}
}

func getEntriesCollection(session *mgo.Session) *mgo.Collection {
	return db.GetMongoCollection(session, FeedEntriesCollection)
}
