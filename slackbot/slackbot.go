package slackbot

import (
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

type Action func(*Bot, *slack.Msg, ...string)

type Bot struct {
	Name   string
	UserID string

	client *slack.Client
	rtm    *slack.RTM
	logger *logrus.Logger

	actions map[*regexp.Regexp]Action
}

func New(token string) *Bot {
	client := slack.New(token)
	logger := logrus.New()

	bot := &Bot{
		client:  client,
		rtm:     client.NewRTM(),
		logger:  logger,
		actions: make(map[*regexp.Regexp]Action),
	}

	return bot
}

func (bot *Bot) Start() {
	bot.handleRTM()
}

func (bot *Bot) RespondTo(match string, action Action) {
	bot.actions[regexp.MustCompile(match)] = action
}

func (bot *Bot) Message(channel string, msg string) {
	bot.client.PostMessage(channel, msg, slack.NewPostMessageParameters())
}

func (bot *Bot) handleRTM() {
	var filter filterer

	rtm := bot.rtm
	log := bot.logger

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			bot.UserID = ev.Info.User.ID
			bot.Name = ev.Info.User.ID

			filter = newDirectFilter(bot.UserID)

			log.Infof("%s is online @ %s", bot.Name, ev.Info.Team.Name)
			log.Debugln("Bot info:", ev.Info)
			log.Debugln("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			if filter.filter(&ev.Msg) {
				log.Debugf("Message: %v\n", ev)
				bot.handleMsg(&ev.Msg)
			}

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

func (bot *Bot) handleMsg(msg *slack.Msg) {
	txt := bot.cleanupMsg(msg.Text)

	for match, action := range bot.actions {
		if matches := match.FindAllStringSubmatch(txt, -1); matches != nil {
			action(bot, msg, matches[0]...)
			return
		}
	}
}

func (bot *Bot) cleanupMsg(msg string) string {
	return strings.TrimLeft(strings.TrimSpace(msg), "<@"+bot.UserID+"> ")
}
