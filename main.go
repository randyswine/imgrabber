package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	// Source of images links and directory of desctination.
	srcurl, dstdir string
)

func init() {
	// init app flag.
	flag.StringVar(&srcurl, "src", "", "this is http address source of images")
	flag.StringVar(&dstdir, "to", "", "this is directory of download destionation")
}

func main() {
	flag.Parse()
	if srcurl == "" || dstdir == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	//
	source, err := request(srcurl)
	if err != nil {
		fmt.Printf("Error download images from %s: %v", srcurl, err)
		os.Exit(2)
	}
	//
	imgLinks, err := findImgLinks(source)
	if err != nil {
		fmt.Printf("Error of download images from %s: %v\r\n", srcurl, err)
		os.Exit(2)
	}
	//
	for _, imgLink := range imgLinks {
		go download(imgLink, dstdir)
	}
	os.Exit(0)
}

// request calls the GET method on the given URL.
func request(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error request to %s:%v", url, err)
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error receive of response: %v", err)
	}
	return content, nil
}

// findImgLinks parses the source, and searches in it for links to images.
func findImgLinks(source []byte) ([]string, error) {
	return nil, nil
}

// download loads the remote file to the specified directory.
func download(url string, dirpath string) {}
