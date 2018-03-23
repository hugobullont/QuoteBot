package main

import (
	"io"
	"os"

	"github.com/jonas747/dca"
)

func main() {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 120
	opts.Volume = 180 //Volumen del audio

	encodeSession, _ := dca.EncodeFile("test.mp3", opts)
	// Make sure everything is cleaned up, that for example the encoding process if any issues happened isnt lingering around
	defer encodeSession.Cleanup()

	output, err := os.Create("test.dca")
	if err != nil {
		// Handle the error
	}

	io.Copy(output, encodeSession)
}
