package main

import (
	"fmt"
	"github.com/asticode/go-astideepspeech"
	"github.com/gorilla/websocket"
	"github.com/zenwerk/go-wave"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{} // use default options

func socketHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrade to websocket conn
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// Initialize DeepSpeech
	dsModel := astideepspeech.New(*model, nCep, nContext, *alphabet, beamWidth)
	defer dsModel.Close()
	if *lm != "" {
		dsModel.EnableDecoderWithLM(*alphabet, *lm, *trie, lmWeight, validWordCountWeight)
	}
	// Setup stream for input
	streamIn := astideepspeech.SetupStream(dsModel, 0, 16000)


	// Setup Wave file writer
	waveFile, _ := os.Create("fromOut.wav")
	param := wave.WriterParam{
		Out:           waveFile,
		Channel:       1,
		SampleRate:    16000,
		BitsPerSample: 16, // if 16, change to WriteSample16()
	}
	waveWriter, _ := wave.NewWriter(param)


	keepRunning := true

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt == websocket.TextMessage {

			if string(message) == "done" {
				keepRunning = false
				text := streamIn.FinishStream()
				fmt.Println(text)
				err = c.WriteMessage(websocket.TextMessage, []byte(text))
				waveWriter.Close()

				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		}

		if mt == websocket.BinaryMessage && keepRunning {
			//stream message into deepspeech

			//convert message to int16[]
			int16Message := convertByteToInt16(message)

			// write to wave file
			_, err := waveWriter.WriteSample16(int16Message) // WriteSample16 for 16 bits
			errCheck(err)

			//write to deepspeech stream
			streamIn.FeedAudioContent(int16Message, uint(len(int16Message)))
		}

	}
	fmt.Println("BROKE OUT")
}