package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
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
	timeline := time.Now()
	flag.Parse()
	if srcurl == "" || dstdir == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	// receive source of images links
	source, err := request(srcurl)
	if err != nil {
		fmt.Printf("Error download images from %s: %v", srcurl, err)
		os.Exit(2)
	}
	// searching of images links in source
	srcReader := bytes.NewReader(source)
	imgLinks, err := findImgLinks(srcReader)
	if len(imgLinks) == 0 {
		fmt.Printf("Imgrabber cannot found images in %s", srcurl)
		os.Exit(0)
	}
	if err != nil {
		fmt.Printf("Error of download images from %s: %v\r\n", srcurl, err)
		os.Exit(2)
	}
	// creates a directory to load
	err = makeDestinationDir(dstdir)
	if err != nil {
		fmt.Printf("Error of download images from %s: %v\r\n", srcurl, err)
		os.Exit(2)
	}
	// creates a abs urls for downloading
	imgLinks, err = makeAbsURLs(srcurl, imgLinks)
	if err != nil {
		fmt.Printf("Error of download images from %s: %v\r\n", srcurl, err)
		os.Exit(2)
	}
	for _, imgLink := range imgLinks {
		err := download(imgLink, dstdir)
		if err != nil {
			fmt.Printf("Error of download images from %s: %v\r\n", srcurl, err)
			os.Exit(2)
		}
		break
	}
	fmt.Printf("download all files by %.2fs", time.Since(timeline).Seconds())
	os.Exit(0)
}

// makeAbsURLs creates a abs ulrs for images downloading.
func makeAbsURLs(rawURL string, links []string) ([]string, error) {
	var absURLs []string
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("error of creates abs url for downloading: %v", err)
	}
	root := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	for _, link := range links {
		if strings.HasPrefix(link, "/") {
			link = strings.Replace(link, "/", "", 1)
		}
		absURLs = append(absURLs, fmt.Sprintf("%s/%s", root, link))
	}
	return absURLs, nil
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
func findImgLinks(srcReader io.Reader) ([]string, error) {
	var imgLinks []string
	var walker func(links []string, n *html.Node) []string
	doc, err := html.Parse(srcReader)
	if err != nil {
		return nil, fmt.Errorf("error of parse sourse: %v", err)
	}
	// TODO: walker cannot parses iframe.
	walker = func(links []string, n *html.Node) []string {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					links = append(links, a.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			links = walker(links, c)
		}
		return links
	}
	imgLinks = walker(imgLinks, doc)
	return imgLinks, nil
}

// download loads the remote file to the specified directory.
func download(url string, dirpath string) error {
	fmt.Println("download", url)
	imgContent, err := request(url)
	if err != nil {
		return fmt.Errorf("error of download images: %v", err)
	}
	splitURL := strings.Split(url, "/")
	fileName := fmt.Sprintf("%s/%s", dirpath, splitURL[len(splitURL)-1])
	err = ioutil.WriteFile(fileName, imgContent, os.ModePerm)
	return nil
}

// makeDestinationDir creates a directory to load.
func makeDestinationDir(dirpath string) error {
	if isPathExists(dirpath) == false {
		err := os.Mkdir(dirpath, os.ModeDir)
		if err != nil {
			return fmt.Errorf("error of download images (make dirpath): %v", err)
		}
	}
	return nil
}

/*
	isPathExists returns a boolean indicating whether the error is
	known to report that a file or directory does not exist.
*/
func isPathExists(path string) bool {
	result := true
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		result = false
	}
	return result
}
