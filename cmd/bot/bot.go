package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

func main() {
	bot := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	if _, yes := os.LookupEnv("SLACK_BOT_DEBUG"); yes {
		log.SetLevel(log.DebugLevel)
		bot.SetDebug(true)
	}

	rtm := bot.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello event

		case *slack.ConnectedEvent:
			log.Infoln("Connected!")
			log.Debugln("Infos:", ev.Info)
			log.Debugln("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			log.Debug("Message: %v\n", ev)

		case *slack.RTMError:
			log.Errorf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Fatalln("Invalid credentials")
			return

		default:
			// Can be used to handle custom events.
			// See: https://github.com/danackerson/bender-slackbot
		}
	}
}
