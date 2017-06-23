package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/develed/develed/config"
	srv "github.com/develed/develed/services"
	"github.com/develed/develed/slackbot"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	cfg = flag.String("config", "/etc/develed.toml", "configuration file")
)

func main() {
	flag.Parse()

	conf, err := config.Load(*cfg)
	if err != nil {
		log.Fatalln(err)
	}

	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		token = conf.Bot.SlackToken
	}

	bot := slackbot.New(token)

	textdConn, err := grpc.Dial(conf.Textd.GRPCServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer textdConn.Close()

	imagedConn, err := grpc.Dial(conf.Imaged.GRPCServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer imagedConn.Close()

	textd := srv.NewTextdClient(textdConn)
	imaged := srv.NewImagedClient(imagedConn)

	bot.DefaultResponse(func(b *slackbot.Bot, msg *slack.Msg) {
		bot.Message(msg.Channel, "Non ho capito")
	})

	bot.RespondTo("scrivi (.*)", func(b *slackbot.Bot, msg *slack.Msg, args ...string) {
		text := args[1]

		_, err := textd.Write(context.Background(), &srv.TextRequest{
			Text: text,
		})
		if err != nil {
			log.Errorln(err)
		} else {
			bot.Message(msg.Channel, "Hai scritto: "+text)
		}
	})

	bot.RespondTo("mostra (http{s?}://.*)", func(b *slackbot.Bot, msg *slack.Msg, args ...string) {
		url := args[1]

		_, err := imaged.Show(context.Background(), &srv.ImageRequest{
			Source: &srv.ImageRequest_Url{Url: url},
		})
		if err != nil {
			log.Errorln(err)
		} else {
			bot.Message(msg.Channel, ":thumbsup:")
		}
	})

	bot.Start()
}
