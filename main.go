package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Repo is a representation of a github repo directory
type Repo struct {
	Name   string
	Route  string
	Commit int64
	Dirs   []Dir
	Files  []string
}

// Dir is a representation of a directory and it's content in a repo
type Dir struct {
	Name  string
	route string
	Dirs  []Dir
	Files []string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	repoName := "Make-School-Courses/SPD-2.31-Testing-and-Architecture"
	baseURL := "https://github.com/" + repoName
	repo := Repo{Name: repoName, Route: "master"}
	initRepo(baseURL, &repo)
	fmt.Println(repo)

}

func initRepo(baseURL string, repo *Repo) {
	c := colly.NewCollector()
	c.OnHTML("a[data-pjax] span strong", func(e *colly.HTMLElement) {
		commit, err := strconv.ParseInt(string(e.Text), 10, 16)
		check(err)
		repo.Commit = commit
	})
	c.OnHTML("span a.js-navigation-open", func(e *colly.HTMLElement) {
		fileType := strings.Split(e.Attr("href"), "/")[3]
		title := e.Attr("title")
		if fileType == "blob" {
			repo.Files = append(repo.Files, title)
		} else {
			repo.Dirs = append(repo.Dirs, Dir{Name: title})
		}

	})
	c.Visit(baseURL)

}

func raw(path string) string {
	return "https://raw.githubusercontent.com/" + path
}
