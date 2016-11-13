package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dankeroni/gotwitch"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/common/config"
	"github.com/pajlada/pajbot2/plog"
	"github.com/pajlada/pajbot2/redismanager"
	"github.com/shogo82148/go-shuffle"
)

var log = plog.GetLogger()

var configPath = flag.String("config", "./config.json", "")

func main() {
	plog.InitLogging()

	flag.Parse()

	runCmd()
}

var quitChannel chan string

// Redis xD
var Redis *redismanager.RedisManager
var shuffler *shuffle.Shuffler

func init() {
	shuffler = shuffle.New(rand.NewSource(int64(time.Now().Nanosecond())))
}

func runCmd() {
	// TODO: Use config path from system arguments
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("An error occured while loading the config file:", err)
	}

	Redis = redismanager.Init(config)

	// Initialize twitch API
	apirequest.Twitch = gotwitch.New(config.Auth.Twitch.User.ClientID)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		quitChannel <- "Quitting due to SIGTERM/SIGINT"
	}()
	quitChannel = make(chan string)

	// Start "services"
	go runTwitch()
	go runBTTV()
	go runFFZ()

	q := <-quitChannel
	log.Fatal(q)
}
