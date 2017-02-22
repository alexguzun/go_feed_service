package repository

import (
	"github.com/alexguzun/go_feed_service/domain/models"
)

type FeedEntriesRepository interface {
	IsNewEntry(feedEntryId string) (bool, error)
	Save(entry models.FeedEntry)
}
