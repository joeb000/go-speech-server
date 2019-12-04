package main

import (
	"flag"
	"fmt"
	"github.com/asticode/go-astideepspeech"
	"html"
	"log"
	"net/http"
)

var configDir = flag.String("configDir", "", "Path to config directory for the DeepSpeech Model files")
var model = flag.String("model", "output_graph.pbmm", "File name of the model (protocol buffer binary file)")
var alphabet = flag.String("alphabet", "alphabet.txt", "File name of the configuration file specifying the alphabet used by the network")
var lm = flag.String("lm", "lm.binary", "File name of the language model binary file")
var trie = flag.String("trie", "trie", "File name of the language model trie file created with native_client/generate_trie")

var runTime = flag.Int("rt", 0, "Length of time in seconds to listen for audio in before processing")

var M *astideepspeech.Model

func errCheck(err error) {

	if err != nil {
		panic(err)
	}
}

func main() {

	configureFlags()

	// Initialize DeepSpeech
	M = astideepspeech.New(*model, nCep, nContext, *alphabet, beamWidth)
	defer M.Close()
	if *lm != "" {
		M.EnableDecoderWithLM(*alphabet, *lm, *trie, lmWeight, validWordCountWeight)
	}

	http.HandleFunc("/rec", ReceiveFile)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Firend ", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/stream", socketHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
