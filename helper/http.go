package helper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Request calls the GET method on the given URL.
func Request(url string) ([]byte, error) {
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

// MakeAbsURLs creates a abs ulrs for images downloading.
func MakeAbsURLs(rawURL string, links []string) ([]string, error) {
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
