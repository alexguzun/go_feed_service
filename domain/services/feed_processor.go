package services

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alexguzun/go_feed_service/domain/models"
	"github.com/alexguzun/go_feed_service/domain/repository"
	"github.com/mmcdole/gofeed"
	"os"
	"time"
)

const defaultReadInterval = "5m"

var ticker = time.NewTicker(time.Minute * 5)

func init() {
	envValue, isDefined := os.LookupEnv("FEED_READ_INTERVAL")
	var interval string
	if isDefined {
		interval = envValue
	} else {
		interval = defaultReadInterval
	}

	readInterval, err := time.ParseDuration(interval)
	if err != nil {
		panic("Invalid read interval. See https://golang.org/pkg/time/#ParseDuration for valid duration definition")
	} else {
		log.Debugf("Read interval set to %v", interval)
		ticker = time.NewTicker(readInterval)
	}
}

func Start(feedSourceRepo repository.FeedSourceRepository, feedEntriesRepo repository.FeedEntriesRepository) {
	var newFeedsChannel = make(chan models.FeedEntry)

	go func() {
		for range ticker.C {
			processFeeds(feedSourceRepo, newFeedsChannel)
		}
	}()

	go store(feedEntriesRepo, newFeedsChannel)
}

func Stop() {
	ticker.Stop()
}

func readRss(feedSource *models.FeedSource) *gofeed.Feed {
	fp := gofeed.NewParser()

	for _, feed_url := range feedSource.RssUrls {
		fieldss := log.Fields{
			"feed_source": feedSource.Description,
			"url":         feed_url,
		}
		log.WithFields(fieldss).Debug("Reading feeds")

		results, err := fp.ParseURL(feed_url)
		if err != nil {
			log.WithError(err).WithFields(fieldss).Error("Failed to read from url")
			return nil
		}
		return results
	}

	return nil
}

func getFeedEntryUrl(feedEntry *gofeed.Item, feedSource *models.FeedSource) string {
	if feedSource.UseFeedGuid {
		return feedEntry.GUID
	} else {
		return feedEntry.Link
	}
}

func getFeedEntries(feedSource models.FeedSource, newFeedsChannel chan models.FeedEntry, feedSourceRepo repository.FeedSourceRepository) {
	//at the end update the last check date of the feed source
	defer feedSourceRepo.UpdateFeedChecked(&feedSource, time.Now())

	feeds := readRss(&feedSource)
	if feeds == nil || len(feeds.Items) == 0 {
		log.WithField("feed_source", feedSource.Description).Warn("No items from rss feed source")
	} else {
		log.WithField("feed_source", feedSource.Description).Debugf("Got %v items", len(feeds.Items))
		for _, feedEntry := range feeds.Items {
			newFeedsChannel <- models.FeedEntry{
				ID:       getFeedEntryUrl(feedEntry, &feedSource),
				Title:    feedEntry.Title,
				Content:  feedEntry.Description,
				AddedOn:  time.Now(),
				SourceID: feedSource.ID,
			}
		}
	}
}

func processFeeds(feedSourceRepo repository.FeedSourceRepository, newFeedsChannel chan models.FeedEntry) {
	//always read from DB because there might be changes
	feedSources, err := feedSourceRepo.GetAllActive()
	if err != nil {
		log.WithError(err).Panic("Failed to connect to db")
		panic(err)
	}

	log.Debugf("Have %v active feed sources to process", len(feedSources))

	for _, feedSource := range feedSources {
		go getFeedEntries(feedSource, newFeedsChannel, feedSourceRepo)
	}
}

func store(feedEntriesRepo repository.FeedEntriesRepository, newFeedsChannel chan models.FeedEntry) {
	for newFeedEntry := range newFeedsChannel {
		isNew, err := feedEntriesRepo.IsNewEntry(newFeedEntry.ID)
		if err != nil {
			log.WithError(err).WithField("feed_entry_id", newFeedEntry.ID).Error("Failed to query")
		} else if isNew {
			log.WithField("feed_entry_id", newFeedEntry.ID).Debug("New entry")
			feedEntriesRepo.Save(newFeedEntry)
		}
	}
}
