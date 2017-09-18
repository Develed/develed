package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	queue "github.com/develed/develed/queue"

	"github.com/develed/develed/imconv"
	srv "github.com/develed/develed/services"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
)

// TermSink redirects any image written to it to the terminal's stdout.
// It requires a Truecolor-capable terminal in order to render images correctly.
type TermSink struct{}

var q queue.Queue

// NewTermSink creates a new TermSink. It returns an error if the sdout of the
// calling process is not a terminal.

func NewTermSink() (*TermSink, error) {
	//goroutine passo coda
	go DrawRoutine(&q)
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return nil, errors.New("TermSink requires stdout to be a terminal")
	}
	return &TermSink{}, nil
}

func (ts *TermSink) Run() error {
	select {}
}

func (ts *TermSink) Draw(ctx context.Context, req *srv.DrawRequest) (*srv.DrawResponse, error) {
	n := queue.Node{req.Priority, req.Timeslot, req.Data}
	q.Push(&n)
	fmt.Println(n.Priority)
	return &srv.DrawResponse{Code: 0, Status: "OK"}, nil
}

func DrawRoutine(q *queue.Queue) {
	for {
		node := q.Pop()
		if node != nil {
			fmt.Println(int64(node.Priority))
			img, _, err := image.Decode(bytes.NewReader(node.Data))
			if err != nil {
				continue
			}
			fmt.Print("\033[s")
			fmt.Print("\n\n")
			sz := img.Bounds().Size()
			for y := 0; y < sz.Y; y++ {
				for x := 0; x < sz.X; x++ {
					col := imconv.NormalizeColor(img.At(x, y))
					if _, err := fmt.Printf("\033[48;2;%d;%d;%dm ", col.R, col.G, col.B); err != nil {
						continue
					}
				}
				if _, err := fmt.Print("\033[0m\n"); err != nil {
					continue
				}
			}
			fmt.Print("\033[u")
			time.Sleep(time.Duration(node.TimeSlot) * time.Millisecond)
		}
	}
}
