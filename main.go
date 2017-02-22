package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alexguzun/go_feed_service/domain/repository"
	"github.com/alexguzun/go_feed_service/domain/services"
	"github.com/alexguzun/go_feed_service/infrastructure/db"
	repoImpl "github.com/alexguzun/go_feed_service/infrastructure/repository"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Info("Starting")
	log.SetLevel(log.DebugLevel)

	var finishUP = make(chan struct{})
	var gracefulStop = make(chan os.Signal)

	signal.Notify(gracefulStop, os.Interrupt, syscall.SIGTERM)
	signal.Notify(gracefulStop, os.Interrupt, syscall.SIGINT)

	go func() {
		// wait for our os signal to stop the app
		// on the graceful stop channel
		sig := <-gracefulStop
		log.WithField("sig", sig).Debug("Caught signal")

		// send message on "finish up" channel to tell the app to
		// gracefully shutdown
		finishUP <- struct{}{}
	}()

	session := db.GetMongoSession()
	defer session.Close()

	var feedSourceRepo repository.FeedSourceRepository = &repoImpl.FeedSourceMongoRepo{Session: session}
	var feedEntriesRepo repository.FeedEntriesRepository = &repoImpl.FeedEntriesMongoRepo{Session: session}

	services.Start(feedSourceRepo, feedEntriesRepo)

	log.Info("Started")

	// wait for finishUP channel write to close the app down
	<-finishUP
	log.Debug("Stopping things, might take 2 seconds")

	services.Stop()

	time.Sleep(2 * time.Second)

	os.Exit(1)
}
