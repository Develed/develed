package main

import (
	"flag"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/develed/develed/config"
	srv "github.com/develed/develed/services"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	debug = flag.Bool("debug", false, "enter debug mode")
	cfg   = flag.String("config", "/etc/develed.toml", "configuration file")
)

var conf *config.Global

type server struct {
	textd srv.TextdClient
}

func (s *server) Show(ctx context.Context, req *srv.TimeFormat) (*srv.TimeResponse, error) {
	if req.Format == "" {
		req.Format = "15:02"
	}

	resp, err := s.textd.Write(ctx, &srv.TextRequest{
		Text: time.Now().Format(req.Format),
		Font: req.Font,
	})
	if err != nil {
		return nil, err
	}

	return &srv.TimeResponse{
		Code:   resp.Code,
		Status: resp.Status,
	}, nil
}

func main() {
	var err error

	flag.Parse()

	conf, err = config.Load(*cfg)
	if err != nil {
		log.Fatalln(err)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	sock, err := net.Listen("tcp", conf.Timed.GRPCServerAddress)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := grpc.Dial(conf.Textd.GRPCServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	s := grpc.NewServer()
	srv.RegisterTimedServer(s, &server{
		textd: srv.NewTextdClient(conn),
	})
	reflection.Register(s)

	if err := s.Serve(sock); err != nil {
		log.Fatalln(err)
	}
}
