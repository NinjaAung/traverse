package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Dir is an embedded struct of folders in a repo
type Dir struct {
	Name  string
	Route string
	Files []string
	Dirs  []Dir
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

func raw(path string) string {
	return "https://raw.githubusercontent.com/" + path
}

func main() {
	repoName := "torvalds/linux"
	baseURL := "https://github.com/" + repoName
	repo := Repo{Name: repoName}
	dir := Dir{Route: "master"}
	start := time.Now()
	initRepo(baseURL, &repo)
	searchFolder(baseURL, &dir)
	repo.Dir = dir
	saveRepo("./test.json", &repo)
	elapsed := time.Since(start)
	fmt.Println(elapsed)

}

func isFileExists(filePath string) error {
	_, err := os.Open(filePath)
	return err
}

func saveRepo(filePath string, repo *Repo) error {
	if isFileExists(filePath) != nil {
		f, err := os.Create(filePath)
		check(err)
		repoJSON, err := json.MarshalIndent([]*Repo{repo}, "", "  ")
		check(err)
		f.Write(repoJSON)
		f.Close()
	} else {
		return isFileExists(filePath)
	}
	updateJSON(filePath, repo)
	return nil
}

func updateJSON(filePath string, repo *Repo) {
	var repos []*Repo
	f, _ := ioutil.ReadFile(filePath)
	r := []*Repo{repo}
	json.Unmarshal(f, &repos)
	if len(repos) == 0 {
		return
	}
	fillerSize := 5 - len(repos)
	filler := make([]*Repo, fillerSize)
	repos = append(repos, filler...)

	for i, repo := range repos {
		if i >= fillerSize {
			break
		}
		if repo.Name == r[0].Name {
			if i > 4 {
				repos = repos[:4]
			} else {
				repos = append(repos[:i], repos[i+1:]...)
			}
			break
		}
	}
	repos = append(r, repos...)
	if fillerSize > 0 {
		repos = repos[:fillerSize]
	}
	if len(repos) >= 5 {
		repos = repos[:5]
	}
	repoJSON, err := json.MarshalIndent(repos, "", "  ")
	check(err)
	json, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	json.Truncate(0)
	json.Write(repoJSON)
	json.Close()
	check(err)
}

func initRepo(baseURL string, repo *Repo) {
	c := colly.NewCollector()
	c.OnHTML("a[data-pjax] span strong", func(e *colly.HTMLElement) {
		commitsStr := e.Text
		if strings.Contains(commitsStr, ",") {
			commitsStr = strings.Join(strings.Split(commitsStr, ","), "")
		}
		commits, err := strconv.ParseInt(commitsStr, 10, 64)
		check(err)
		repo.Commits = commits
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
