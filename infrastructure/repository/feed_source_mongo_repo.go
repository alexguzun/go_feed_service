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
	var feedSources = []models.FeedSource{}
	err := getSourceCollection(repo).
		Find(bson.M{"active": true}).
		All(&feedSources)

	return feedSources, err
}

func (repo *FeedSourceMongoRepo) UpdateFeedChecked(feedSource *models.FeedSource, whenChecked time.Time) {
	err := getSourceCollection(repo).
		Update(bson.M{"_id": feedSource.ID}, bson.M{"$set": bson.M{"last_check": whenChecked}})
	if err != nil {
		log.WithField("feed_source_id", feedSource.ID).WithError(err).Error("Failed to update feed source")
	}
}

func getSourceCollection(repo *FeedSourceMongoRepo) *mgo.Collection {
	return db.GetMongoCollection(repo.Session, FeedSourceCollection)
}
