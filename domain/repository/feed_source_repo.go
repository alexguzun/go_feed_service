package repository

import (
	"github.com/alexguzun/go_feed_service/domain/models"
	"time"
)

type FeedSourceRepository interface {
	GetAllActive() ([]models.FeedSource, error)
	UpdateFeedChecked(feedSource *models.FeedSource, whenChecked time.Time)
}
