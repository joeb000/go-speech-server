package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/asticode/go-astideepspeech"
	"github.com/asticode/go-astilog"
	"github.com/cryptix/wav"
	"github.com/pkg/errors"
)

// Constants
const (
	beamWidth            = 500
	nCep                 = 26
	nContext             = 9
	lmWeight             = 0.75
	validWordCountWeight = 1.85
)

func metadataToString(m *astideepspeech.Metadata) string {
	retval := ""
	for _, item := range m.Items() {
		retval += item.Character()
	}
	return retval
}

func ReceiveFile(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Processing new request...")
	var Buf bytes.Buffer
	//create a temp file for this
	tmpF, err := ioutil.TempFile("", "audio")
	if err != nil {
		panic(err)
	}

	file, header, err := req.FormFile("file")
	fmt.Printf("\nfile t= %T -- header.Size = %T\n", tmpF, header.Size)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(&Buf, file)
	ioutil.WriteFile(tmpF.Name(), Buf.Bytes(), 0644)

	//splice in code here
	// Stat audio
	ii := header.Size

	// Create reader
	r, err := wav.NewReader(tmpF, ii)
	if err != nil {
		astilog.Fatal(errors.Wrap(err, "creating new reader failed"))
	}

	// Read
	var d []int16
	for {
		// Read sample
		s, err := r.ReadSample()
		if err == io.EOF {
			break
		} else if err != nil {
			astilog.Fatal(errors.Wrap(err, "reading sample failed"))
		}

		// Append
		d = append(d, int16(s))
	}

	output := ""
	// Speech to text

	output = M.SpeechToText(d, uint(len(d)), 44100)

	astilog.Infof("Text: %s", output)

	Buf.Reset()

	io.WriteString(w, fmt.Sprintf("\nText: %s\n\n", output))

}
