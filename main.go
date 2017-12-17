package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"strings"
)

var opts struct {
	Code string `short:"c" long:"code" description:"NGDP Game Code from Wowdev.wiki" required:"true"`
	File string `short:"f" long:"file" description:"CDNS, BGDL, Versions" required:"false" default:"versions"`
}

func main() {
	_, err := flags.Parse(&opts)

	if err != nil || (!isValidFile(opts.File)) {
		fmt.Println("Allowed types: bgdl, versions, cdns")
		return
	}

	raw, _, _ := Version(opts.Code, opts.File)

	fmt.Println(raw)
}

func isValidFile(file string) bool {
	switch strings.ToLower(file) {
	case
		"bgdl",
		"versions",
		"cdns":
		return true
	}
	return false
}
