package main

import (
	"sync"
	"time"

	"github.com/dankeroni/gotwitch"
	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/models"
)

func runTwitch() {
	log.Info("Initializing Twitch services")

	go twitchUpdateStreamStatus()
	// twitchRunStreamStatusUpdate()
}

func twitchRunStreamStatusUpdate() {
	channelList := []string{
		"pajlada",
		"forsenlol",
		"imaqtpie",
		"destiny",
		"tsm_dyrus",
		"neviilz",
		"exbc",
		"pajbot",
	}
	now := time.Now()
	conn := models.GetRedisConn()
	var wg sync.WaitGroup
	defer func(conn redis.Conn) {
		conn.Flush()
		conn.Close()
	}(conn)
	wg.Add(len(channelList))
	for _, channel := range channelList {
		go func(conn redis.Conn, channel string) {
			defer wg.Done()
			apirequest.Twitch.GetStream(channel,
				func(twitchStream gotwitch.Stream) {
					if twitchStream.ID == 0 {
						lastStream, err := models.GetLastStream(channel)
						if err != nil {
							// log.Errorf("Error getting last stream for %s: %s", channel, err)
							return
						}

						if lastStream.End != nil {
							// Last stream is already marked as offline
							return
						}

						lastStream.End = &now
						log.Debugf("%s is now marked as offline", channel)

						models.UpdateLastStreamP(conn, channel, lastStream)
					} else {
						pajbotStreamChunk := models.ChunkFromTwitchStream(twitchStream)
						channel := twitchStream.Channel.Name

						isNew, prevStream, err := models.IsNewStreamChunk(channel, pajbotStreamChunk)
						if err != nil {
							log.Error(err)
							return
						}
						if isNew {
							// Here we add logic to check if last stream ended within the last 5 minutes
							log.Debugf("%#v", prevStream)
							if prevStream != nil {
								const threshold = 10 * time.Minute
								// If the previous stream did not have time to go offline fully, or if 10 minutes passed
								if prevStream.End == nil || prevStream.End.Add(threshold).After(now) {
									prevStream.End = nil
									prevStream.AddStreamChunk(pajbotStreamChunk)
									// Previous stream ended within the last 5 minutes, just continue on that one instead
									log.Debugf("[%s] Updating last stream with new chunk %d", channel, pajbotStreamChunk.ID)
									models.UpdateLastStreamP(conn, channel, prevStream)
									return
								}
							}

							// Add new stream with new chunk
							log.Debugf("Adding new stream with new chunk %s - %d", channel, pajbotStreamChunk.ID)
							models.AddNewStreamP(conn, channel, pajbotStreamChunk)
						} else {
							// log.Debugf("[%s] Update last seen for chunk %d", channel, pajbotStreamChunk.ID)
							prevStream.End = nil
							prevStream.UpdateStreamChunk(pajbotStreamChunk)
							models.UpdateLastStreamP(conn, channel, prevStream)
						}
					}
				},
				func(statusCode int, statusMessage, errorMessage string) {
					log.Debugf("ERROR: %d", statusCode)
					log.Debug(statusMessage)
					log.Debug(errorMessage)
				}, func(err error) {
					log.Debug("Internal error")
				})
		}(conn, channel)
	}

	wg.Wait()
}

func twitchUpdateStreamStatus() {
	const streamStatusInterval = 1 * time.Minute
	ticker := time.NewTicker(streamStatusInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				// Figure out which streams to poll
				// Let's poll at most 10 at a time
				twitchRunStreamStatusUpdate()
			}
		}
	}()
}
