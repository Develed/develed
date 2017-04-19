package main

import (
	"os"

	"github.com/nlopes/slack"
	"github.com/plorefice/develed/slackbot"
)

func main() {
	bot := slackbot.New(os.Getenv("SLACK_BOT_TOKEN"))

	bot.RespondTo("ciao", func(b *slackbot.Bot, msg *slack.Msg, args ...string) {
		bot.Message(msg.Channel, "Olà!")
	})

	bot.RespondTo("come va?", func(b *slackbot.Bot, msg *slack.Msg, args ...string) {
		bot.Message(msg.Channel, "Tutto bene, tu?")
	})

	bot.Start()
}
