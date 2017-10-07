package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"
	"github.com/develed/develed/config"
	srv "github.com/develed/develed/services"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	cfg   = flag.String("config", "/etc/develed.toml", "configuration file")
	debug = flag.Bool("debug", false, "Write to textd")
)

type server struct {
	sink srv.ImageSinkClient
}

func (s *server) Write(ctx context.Context, req *srv.TextRequest) (*srv.TextResponse, error) {
	return &srv.TextResponse{
		Code:   0,
		Status: "Ok",
	}, nil
}

func main() {
	flag.Parse()
	conf, err := config.Load(*cfg)
	if err != nil {
		log.Fatalln(err)
	}

	text := "ciao"
	if flag.NArg() >= 1 {
		text = flag.Arg(0)
	}

	conn, err := grpc.Dial(conf.Textd.GRPCServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	textd := srv.NewTextdClient(conn)
	resp, err := textd.Write(context.Background(), &srv.TextRequest{
		Text:      text,
		FontColor: 0xFFAABBCC,
		FontBg:    0x00112233,
	})
	if err != nil {
		log.Errorln(err)
	}
	log.Info(resp.Status, " ", resp.Code)
	//
	//conn, err := grpc.Dial(conf.DSPD.GRPCServerAddress, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer conn.Close()

	//s := grpc.NewServer()
	//drawing_srv := &server{sink: srv.NewImageSinkClient(conn)}
	//srv.RegisterTextdServer(s, drawing_srv)
	//reflection.Register(s)

	//frame := image.NewRGBA(image.Rect(0, 0, 39, 9))
	//draw.Draw(frame, frame.Bounds(), &image.Uniform{color.RGBA{0, 255, 0, 255}}, image.ZP, draw.Src)
	//buf := &bytes.Buffer{}
	//png.Encode(buf, frame)
	//resp, err := drawing_srv.sink.Draw(context.Background(), &srv.DrawRequest{
	//	Priority: int64(10),
	//	Timeslot: int64(9),
	//	Data:     buf.Bytes(),
	//})

	//frame = image.NewRGBA(image.Rect(0, 0, 39, 9))
	//draw.Draw(frame, frame.Bounds(), &image.Uniform{color.RGBA{255, 0, 255, 255}}, image.ZP, draw.Src)
	//buf = &bytes.Buffer{}
	//png.Encode(buf, frame)
	//resp, err = drawing_srv.sink.Draw(context.Background(), &srv.DrawRequest{
	//	Priority: int64(10),
	//	Timeslot: int64(9),
	//	Data:     buf.Bytes(),
	//})

	//frame = image.NewRGBA(image.Rect(0, 0, 39, 9))
	//draw.Draw(frame, frame.Bounds(), &image.Uniform{color.RGBA{200, 50, 180, 255}}, image.ZP, draw.Src)
	//buf = &bytes.Buffer{}
	//png.Encode(buf, frame)
	//resp, err = drawing_srv.sink.Draw(context.Background(), &srv.DrawRequest{
	//	Priority: int64(10),
	//	Timeslot: int64(9),
	//	Data:     buf.Bytes(),
	//})

	//frame = image.NewRGBA(image.Rect(0, 0, 39, 9))
	//draw.Draw(frame, frame.Bounds(), &image.Uniform{color.RGBA{10, 10, 10, 255}}, image.ZP, draw.Src)
	//buf = &bytes.Buffer{}
	//png.Encode(buf, frame)
	//resp, err = drawing_srv.sink.Draw(context.Background(), &srv.DrawRequest{
	//	Priority: int64(10),
	//	Timeslot: int64(9),
	//	Data:     buf.Bytes(),
	//})
	//fmt.Println(text, resp.Status)

}
