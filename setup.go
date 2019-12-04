package main

import (
	"flag"
	"fmt"
	"os"
)

func addConfigPath(fileName *string) {
	*fileName = fmt.Sprintf("%v/%v", *configDir, *fileName)
}

func configureFlags() {
	flag.Parse()

	if *configDir == "" {
		fmt.Println("No config Directory Specified")
		flag.PrintDefaults()
		os.Exit(0)
	}

	addConfigPath(model)
	addConfigPath(alphabet)
	addConfigPath(lm)
	addConfigPath(trie)
}