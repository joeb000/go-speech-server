package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/asticode/go-astideepspeech"
	"github.com/asticode/go-astilog"
)

var model = flag.String("model", "", "Path to the model (protocol buffer binary file)")
var alphabet = flag.String("alphabet", "", "Path to the configuration file specifying the alphabet used by the network")
var audio = flag.String("audio", "", "Path to the audio file to run (WAV format)")
var lm = flag.String("lm", "", "Path to the language model binary file")
var trie = flag.String("trie", "", "Path to the language model trie file created with native_client/generate_trie")
var version = flag.Bool("version", false, "Print version and exits")
var extended = flag.Bool("extended", false, "Use extended metadata")

var M *astideepspeech.Model

func main() {
	flag.Parse()

	astilog.FlagInit()

	if *version {
		astideepspeech.PrintVersions()
		return
	}

	if *model == "" || *alphabet == "" {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
