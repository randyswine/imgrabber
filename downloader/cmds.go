package downloader

type Cmd int8

const (
	_       = iota
	RUN Cmd = +1
	STOP
)
