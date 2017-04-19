package main

import (
	"os"

	"github.com/develed/develed/slackbot"
	"github.com/nlopes/slack"
)

func main() {
	bot := slackbot.New(os.Getenv("SLACK_BOT_TOKEN"))

	bot.DefaultResponse(func(b *slackbot.Bot, msg *slack.Msg) {
		bot.Message(msg.Channel, "Non ho capito")
	})

	bot.RespondTo("ciao", func(b *slackbot.Bot, msg *slack.Msg, args ...string) {
		bot.Message(msg.Channel, "Olà!")
	})

	bot.RespondTo("echo (.*)", func(b *slackbot.Bot, msg *slack.Msg, args ...string) {
		bot.Message(msg.Channel, "Hai scritto: "+args[1])
	})

	bot.Start()
}
