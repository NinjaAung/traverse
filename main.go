package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Dir is an embedded struct of folders in a repo
type Dir struct {
	Name  string
	Route string
	Dirs  []Dir
	Files []string
}

// Repo is a representation of a github repo directory
type Repo struct {
	Name    string
	Commits int64
	Dir
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	repoName := "Make-School-Courses/SPD-2.31-Testing-and-Architecture"
	baseURL := "https://github.com/" + repoName
	repo := Repo{Name: repoName}
	dir := Dir{}
	initRepo(baseURL, &repo)
	searchFolder(baseURL, &dir)
	repo.Dir = dir
	fmt.Println(repo)

}

func initRepo(baseURL string, repo *Repo) {
	c := colly.NewCollector()
	c.OnHTML("a[data-pjax] span strong", func(e *colly.HTMLElement) {
		commit, err := strconv.ParseInt(string(e.Text), 10, 16)
		check(err)
		repo.Commits = commit
	})
	c.Visit(baseURL)

}

func searchFolder(link string, dir *Dir) {
	c := colly.NewCollector()
	c.OnHTML("span a.js-navigation-open", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		link := strings.Split(href, "/")
		dir.Route = strings.Join(link[4:len(link)-1], "/")
		routeType := strings.Split(href, "/")[3]
		title := e.Attr("title")
		if routeType == "blob" {
			dir.Files = append(dir.Files, title)
		} else if routeType == "tree" {
			newDir := Dir{Name: title}
			searchFolder("https://github.com"+href, &newDir)
			dir.Dirs = append(dir.Dirs, newDir)

		}

	})
	c.Visit(link)

}

func raw(path string) string {
	return "https://raw.githubusercontent.com/" + path
}
