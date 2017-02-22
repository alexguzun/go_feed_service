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
	count, err := getEntriesCollection(repo).
		Find(bson.M{"_id": feedEntryId}).
		Count()

	if err != nil {
		return false, err
	} else if count == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (repo *FeedEntriesMongoRepo) Save(newFeedEntry models.FeedEntry) {
	err := getEntriesCollection(repo).Insert(newFeedEntry)
	if err != nil {
		log.WithField("feed_entry_id", newFeedEntry.ID).WithError(err).Error("Failed to add new feed entry")
	}
}

func getEntriesCollection(repo *FeedEntriesMongoRepo) *mgo.Collection {
	return db.GetMongoCollection(repo.Session, FeedEntriesCollection)
}
