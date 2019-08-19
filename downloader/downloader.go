package downloader

import (
	"fmt"
	"imgrabber/helper"
	"io/ioutil"
	"os"
	"strings"
)

type downloader struct {
	cmds           <-chan Cmd
	tasks          <-chan string
	Result         chan<- string
	isReceiveLinks bool
	dirpath        string
}

func New(cmds <-chan Cmd, tasks <-chan string) *downloader {
	return &downloader{cmds: cmds,
		tasks:          tasks,
		Result:         make(chan<- string),
		isReceiveLinks: false}
}

func (d *downloader) Listen() {
	go func() {
		for {
			cmd, ok := <-d.cmds
			if !ok {
				return
			}
			if cmd == STOP {
				d.isReceiveLinks = false
				return
			}
			if cmd == RUN {
				d.isReceiveLinks = true
			}
			if d.isReceiveLinks {
				link, ok := <-d.tasks
				if !ok {
					return
				}
				d.download(link)
			}
		}
	}()
}

// download loads the remote file to the specified directory.
func (d *downloader) download(url string) {
	imgContent, err := helper.Request(url)
	if err != nil {
		d.Result <- fmt.Sprintf("error of download images: %v", err)
		return
	}
	splitURL := strings.Split(url, "/")
	fileName := fmt.Sprintf("%s/%s", d.dirpath, splitURL[len(splitURL)-1])
	err = ioutil.WriteFile(fileName, imgContent, os.ModePerm)
	d.Result <- fmt.Sprintf("downloaded %s", url)
	return
}
