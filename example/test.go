package main

import (
	"encoding/gob"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
)

func main() {
	fmt.Println("Run..")

	// Actually read-only, write flag required to avoid blocking on open()
	fifo, err := os.OpenFile(os.Args[1], os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Fatalln(err)
	}
	defer fifo.Close()

	enc := gob.NewEncoder(fifo)
	m, _, err := image.Decode(os.Stdin)
	if err != nil {
		panic(err)
	}

	gob.Register(&image.RGBA{})
	enc.Encode(&m)
}
