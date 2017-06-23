package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/develed/develed/config"
	srv "github.com/develed/develed/services"
	"github.com/nfnt/resize"
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
	sink srv.ImageSinkClient
}

func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to retrieve image with code %d\n", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (s *server) Show(ctx context.Context, req *srv.ImageRequest) (*srv.ImageResponse, error) {
	var err error
	var buf []byte

	switch t := req.Source.(type) {
	case *srv.ImageRequest_Data:
		buf = req.GetData()
	case *srv.ImageRequest_Url:
		buf, err = downloadImage(req.GetUrl())
		if err != nil {
			return nil, err
		}
	case nil:
		log.Warnln("Empty request!")
		return nil, errors.New("Invalid empty request")
	default:
		log.Errorf("Image source has unexpected type %T\n", t)
		return nil, errors.New("Invalid type in request")
	}

	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	png.Encode(&out, resize.Resize(39, 9, img, resize.Lanczos3))

	resp, err := s.sink.Draw(context.Background(), &srv.DrawRequest{
		Data: out.Bytes(),
	})
	if err != nil {
		return nil, err
	}

	return &srv.ImageResponse{
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

	sock, err := net.Listen("tcp", conf.Imaged.GRPCServerAddress)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := grpc.Dial(conf.DSPD.GRPCServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	s := grpc.NewServer()
	srv.RegisterImagedServer(s, &server{
		sink: srv.NewImageSinkClient(conn),
	})
	reflection.Register(s)

	if err := s.Serve(sock); err != nil {
		log.Fatalln(err)
	}
}
