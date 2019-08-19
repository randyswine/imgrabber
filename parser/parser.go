package parser

import (
	"bytes"
	"fmt"
	"imgrabber/helper"

	"golang.org/x/net/html"
)

type parser struct {
	tagLink   string
	srcReader *bytes.Reader
}

func New(tagLink string, source []byte) *parser {
	return &parser{tagLink: tagLink, srcReader: bytes.NewReader(source)}
}

// FindImgLinks parses the source, and searches in it for links to images.
func (p *parser) FindImgLinks() ([]string, error) {
	var imgLinks []string
	var walker func(links []string, n *html.Node) []string
	doc, err := html.Parse(p.srcReader)
	if err != nil {
		return nil, fmt.Errorf("error of parse sourse: %v", err)
	}
	//TODO:
	/*
		1. walker cannot parses iframe.
		2. walker not safe use tagLine param.
	*/
	walker = func(links []string, n *html.Node) []string {
		if n.Type == html.ElementNode && n.Data == p.tagLink {
			for _, a := range n.Attr {
				if a.Key == "href" || a.Key == "src" {
					if helper.IsCorrectExtension(a.Val, ".jpg", ".png", ".bmp") {
						links = append(links, a.Val)
					}
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
