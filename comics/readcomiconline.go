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

type ReadComicOnline struct {
	issueMap   map[string]string
	handleFunc func(model.Episode, error)
	base       string
}

func (rco *ReadComicOnline) SetHandler(handleFunc func(episodes model.Episode, err error)) {
	rco.handleFunc = handleFunc
}

func (rco *ReadComicOnline) FindAllIssues(listURL string) error {
	list, err := url.Parse(listURL)
	if err != nil {
		return err
	}
	rco.base = list.Scheme + "://" + list.Host
	fmt.Println(rco.base)

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{listURL},
		ParseFunc: rco.parseFindAllIssues,
	}).Start()

	return nil
}

func (rco *ReadComicOnline) parseFindAllIssues(g *geziyor.Geziyor, r *client.Response) {
	if rco.issueMap == nil {
		rco.issueMap = make(map[string]string)
	}

	r.HTMLDoc.Find("ul>li>div").Each(
		func(i int, s *goquery.Selection) {
			issueLink, _ := s.Find("a").Attr("href")
			title, _ := s.Find("img").Attr("alt")
			rco.issueMap[rco.base+issueLink+"&readType=1"] = fmt.Sprintf("[%d]", i+1) + title
		},
	)
}

func (rco *ReadComicOnline) DownloadIssues() {
	for issueURL := range rco.issueMap {
		fmt.Println("downloading...", issueURL)
		geziyor.NewGeziyor(&geziyor.Options{
			StartURLs: []string{issueURL},
			StartRequestsFunc: func(g *geziyor.Geziyor) {
				g.GetRendered(issueURL, g.Opt.ParseFunc)
			},
			ParseFunc: rco.parseIssue,
		}).Start()

		return
	}
}

func (rco *ReadComicOnline) parseIssue(g *geziyor.Geziyor, r *client.Response) {
	var err error

	totalImages := len(r.HTMLDoc.Find("#divImage").Nodes)
	r.HTMLDoc.Find("#divImage").Find("p>img").Each(
		func(i int, s *goquery.Selection) {
			src, ok := s.Attr("src")
			if !ok {
				return
			}
			fmt.Println(src)

			var u *url.URL
			u, err = url.Parse(r.JoinURL(src))
			if err != nil {
				rco.handleFunc(model.Episode{}, err)
				return
			}

			fmt.Println("downloading...", u)
			req := &http.Request{
				Method: "GET",
				URL:    u,
			}

			resp, err := g.Client.Do(req)
			if err != nil {
				rco.handleFunc(model.Episode{}, err)
				return
			}

			episode := model.NewEpisode(resp, totalImages, err)
			rco.handleFunc(episode, err)
		},
	)
}
