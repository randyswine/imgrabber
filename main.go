package main

import (
	"flag"
	"os"
)

var (
	// Source of images links and directory of desctination.
	srcurl, dstdir, tagLink string
)

func init() {
	// init app flag.
	flag.StringVar(&srcurl, "src", "", "this is http address source of images")
	flag.StringVar(&dstdir, "to", "", "this is directory of download destionation")
	flag.StringVar(&tagLink, "tag", "", "this is the tag in which the image link is defined")
}

func main() {
	flag.Parse()
	if srcurl == "" || dstdir == "" || tagLink == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	os.Exit(0)
}
