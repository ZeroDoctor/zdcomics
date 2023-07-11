package main

import (
	"fmt"

	"github.com/zerodoctor/zdcomics/comics"
	"github.com/zerodoctor/zdcomics/model"
)

type Comics interface {
	SetHandler(handleFunc func(episodes model.Episode, err error))
	FindAllIssues(listURL string)
	DownloadIssues()
}

func main() {
	fmt.Println("vim-go")

	comic := comics.ReadComicOnline{}
	comic.SetHandler(func(episode model.Episode, err error) {
		fmt.Println("found episode...", err)
	})

	issueList := "https://readcomiconline.li/Comic/Invincible"
	comic.FindAllIssues(issueList)
	comic.DownloadIssues()
}
