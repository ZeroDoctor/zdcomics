package comics

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/zerodoctor/zdcomics/model"
)

type Webtoons struct {
	issueMap   map[string]string
	handleFunc func(model.Episode, error)
}

func (w *Webtoons) SetHandler(handleFunc func(episodes model.Episode, err error)) {
	w.handleFunc = handleFunc
}

func (w *Webtoons) FindAllIssues(listURL string) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{listURL},
		ParseFunc: w.parseFindAllIssues,
	}).Start()
}

func (w *Webtoons) parseFindAllIssues(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("#topEpisodeList").Find("div.episode_cont").Find("li").Each(
		func(i int, s *goquery.Selection) {
			num := i + 1
			issueLink, _ := s.Find("a").Attr("href")
			title, _ := s.Find("img").Attr("alt")
			w.issueMap[issueLink] = fmt.Sprintf("[%d]", num) + title
		},
	)
}

func (w *Webtoons) DownloadIssues() {
	for issueURL := range w.issueMap {
		geziyor.NewGeziyor(&geziyor.Options{
			StartURLs: []string{issueURL},
			ParseFunc: w.parseIssue,
		}).Start()
	}
}

func (w *Webtoons) parseIssue(g *geziyor.Geziyor, r *client.Response) {
	var err error

	totalImages := len(r.HTMLDoc.Find("#_imageList").Find("img").Nodes)
	r.HTMLDoc.Find("#_imageList").Find("img").Each(
		func(i int, s *goquery.Selection) {
			href, ok := s.Attr("data-url")
			if !ok {
				return
			}

			var u *url.URL
			u, err = url.Parse(r.JoinURL(href))
			if err != nil {
				return
			}

			req := &http.Request{
				Method: "GET",
				Header: http.Header(map[string][]string{
					"Referer": {"http://www.webtoons.com"}, // ! important header
				}),
				URL: u,
			}

			resp, err := g.Client.Do(req)
			if err != nil {
				return
			}

			episode := model.NewEpisode(resp, totalImages, err)
			w.handleFunc(episode, err)
		},
	)
}
