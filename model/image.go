package model

import "net/url"

type Image struct {
	Url    *url.URL
	Image  []byte
	Type   string
	Width  float64
	Height float64
}
