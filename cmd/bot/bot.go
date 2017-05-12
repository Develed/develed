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

	bot := slackbot.New(os.Getenv("SLACK_BOT_TOKEN"))

	conn, err := grpc.Dial(conf.Textd.GRPCServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	textd := srv.NewTextdClient(conn)

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

	bot.Start()
}
