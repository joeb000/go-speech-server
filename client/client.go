package main

/*
  #include <stdio.h>
  #include <unistd.h>
  #include <termios.h>
  char getch(){
      char ch = 0;
      struct termios old = {0};
      fflush(stdout);
      if( tcgetattr(0, &old) < 0 ) perror("tcsetattr()");
      old.c_lflag &= ~ICANON;
      old.c_lflag &= ~ECHO;
      old.c_cc[VMIN] = 1;
      old.c_cc[VTIME] = 0;
      if( tcsetattr(0, TCSANOW, &old) < 0 ) perror("tcsetattr ICANON");
      if( read(0, &ch,1) < 0 ) perror("read()");
      old.c_lflag |= ICANON;
      old.c_lflag |= ECHO;
      if(tcsetattr(0, TCSADRAIN, &old) < 0) perror("tcsetattr ~ICANON");
      return ch;
  }
*/
import "C"

// stackoverflow.com/questions/14094190/golang-function-similar-to-getchar

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
"flag"
	"github.com/gordonklaus/portaudio"
	wave "github.com/zenwerk/go-wave"
)

var remoteHost = flag.String("host", "localhost:8080", "remote host")


func errCheck(err error) {

	if err != nil {
		panic(err)
	}
}

func convertInt16ToByte(a []int16) []byte {
	b := make([]byte, 2*len(a))
	bI := 0
	for i := 0; i < len(a); i++ {
		b[bI] = byte(uint(a[i]))
		b[bI+1] = byte(uint(a[i] >> 8))
		bI += 2
	}
	return b
}

func main() {
	flag.Parse()


	audioFileName := "outfile"

	fmt.Println("Recording. Press ESC to quit.")

	if !strings.HasSuffix(audioFileName, ".wav") {
		audioFileName += ".wav"
	}
	waveFile, err := os.Create(audioFileName)
	errCheck(err)

	// www.people.csail.mit.edu/hubert/pyaudio/  - under the Record tab
	inputChannels := 1
	outputChannels := 0
	sampleRate := 16000
	framesPerBuffer := make([]int16, 64)

	// init PortAudio

	portaudio.Initialize()
	//defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
	errCheck(err)
	//defer stream.Close()

	// setup Wave file writer

	param := wave.WriterParam{
		Out:           waveFile,
		Channel:       inputChannels,
		SampleRate:    sampleRate,
		BitsPerSample: 16, // if 16, change to WriteSample16()
	}

	waveWriter, err := wave.NewWriter(param)
	errCheck(err)

	// start websocket conn
	u := url.URL{Scheme: "ws", Host: *remoteHost, Path: "/stream"}

	wscon, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer wscon.Close()

	done := make(chan struct{})


	var wg sync.WaitGroup
	wg.Add(1)
	closeSocket := false
	go func() {
		defer close(done)
		for {
			_, message, err := wscon.ReadMessage()

			if err != nil {
				log.Println("read:", err)
				return
			}

			closeSocket = true

			fmt.Printf("\n\nText: %s\n", message)
			err = wscon.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			errCheck(err)
			wg.Done()
		}
	}()

	go func() {
		key := 0

		for key != 27 {
			// better to control
			// how we close then relying on defer
			key = int(C.getch())
			if key == 27 {
				break
			}
			fmt.Println()
			fmt.Println("Cleaning up ...")
		}

		closeSocket = true
		// close ws
		err := wscon.WriteMessage(websocket.TextMessage, []byte("done"))
		errCheck(err)



	}()

	//recording in progress ticker. From good old DOS days.
	ticker := []string{
		"-",
		"\\",
		"/",
		"|",
	}
	rand.Seed(time.Now().UnixNano())

	// start reading from microphone
	errCheck(stream.Start())
	for {
		if closeSocket {
			break
		}
		errCheck(stream.Read())

		fmt.Printf("\rRecording is live now. Say something to your microphone! [%v]", ticker[rand.Intn(len(ticker)-1)])

		// write to wave file
		_, err := waveWriter.WriteSample16(framesPerBuffer) // WriteSample16 for 16 bits
		errCheck(err)

		sBytes := convertInt16ToByte(framesPerBuffer)
		err = wscon.WriteMessage(websocket.BinaryMessage, sBytes)
		errCheck(err)
	}
	fmt.Println("\nDONE\n")
	errCheck(stream.Stop())

	wg.Wait()


	waveWriter.Close()
	stream.Close()
	portaudio.Terminate()
	//fmt.Println("Play", audioFileName, "with a audio player to hear the result.")
	os.Exit(0)

}
