package repository

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alexguzun/go_feed_service/domain/models"
	"github.com/alexguzun/go_feed_service/infrastructure/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type FeedSourceMongoRepo struct {
	Session *mgo.Session
}

const FeedSourceCollection string = "feed_sources"

func (repo *FeedSourceMongoRepo) GetAllActive() ([]models.FeedSource, error) {
	feedSources, err := db.Query(repo.Session, func(session *mgo.Session) (interface{}, error) {
		var results = []models.FeedSource{}
		err := getSourceCollection(session).
			Find(bson.M{"active": true}).
			All(&results)
		return results, err
	})

	return feedSources.([]models.FeedSource), err
}

func (repo *FeedSourceMongoRepo) UpdateFeedChecked(feedSource *models.FeedSource, whenChecked time.Time) {
	err := db.Execute(repo.Session, func(session *mgo.Session) (err error) {
		return getSourceCollection(session).
			Update(bson.M{"_id": feedSource.ID}, bson.M{"$set": bson.M{"last_check": whenChecked}})
	})
	if err != nil {
		log.WithField("feed_source_id", feedSource.ID).WithError(err).Error("Failed to update feed source")
	}
}

func getSourceCollection(session *mgo.Session) *mgo.Collection {
	return db.GetMongoCollection(session, FeedSourceCollection)
}
