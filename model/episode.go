package model

import "net/http"

type Episode struct {
	resp        *http.Response
	totalImages int
	err         error
}

func NewEpisode(resp *http.Response, total int, err error) Episode {
	return Episode{
		resp:        resp,
		totalImages: total,
		err:         err,
	}
}
