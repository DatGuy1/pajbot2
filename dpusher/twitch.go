package main

import (
	"time"

	"github.com/dankeroni/gotwitch"
	"github.com/pajlada/pajbot2/apirequest"
)

func runTwitch() {
	log.Info("Initializing Twitch services")

	go twitchUpdateStreamStatus()
}

func twitchUpdateStreamStatus() {
	const streamStatusInterval = 20 * time.Second
	ticker := time.NewTicker(streamStatusInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				streams := "pajlada,forsenlol,imaqtpie"
				apirequest.Twitch.GetStreams(streams,
					func(streams []gotwitch.Stream) {
						// TODO: Go through streams here and update redis data
						log.Debugf("Num streams: %d", len(streams))
					},
					func(statusCode int, statusMessage, errorMessage string) {
						log.Debugf("ERROR: %d", statusCode)
						log.Debug(statusMessage)
						log.Debug(errorMessage)
					}, func(err error) {
						log.Debug("Internal error")
					})
			}
		}
	}()
}
