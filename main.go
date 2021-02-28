package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	repoName := "Make-School-Courses/SPD-2.31-Testing-and-Architecture"
	baseURL := "https://github.com/" + repoName
	repo := Repo{Name: repoName}
	dir := Dir{Route: "master"}
	start := time.Now()
	initRepo(baseURL, &repo)
	searchFolder(baseURL, &dir)
	repo.Dir = dir
	elapsed := time.Since(start)
	fmt.Println(elapsed)

}

func isFileExists(filePath string) bool {
	_, err := os.Open(filePath)
	return err == nil
}

func saveRepo(name, location string, repo *Repo) {
	var repos []*Repo
	filePath := filepath.Join(location, name)
	if !isFileExists(filePath) {
		f, err := os.Create(filePath)
		check(err)
		repoJSON, err := json.MarshalIndent([]*Repo{repo}, "", "  ")
		check(err)
		f.Write(repoJSON)
		f.Close()
		return
	}
	updateJSON(filePath, repos, repo)
}

func updateJSON(filePath string, repos []*Repo, repo *Repo) {
	f, _ := ioutil.ReadFile(filePath)
	r := []*Repo{repo}
	json.Unmarshal(f, &repos)
	for i, repo := range repos {
		if repo.Name == repo.Name {
			if i > 4 {
				repos = repos[:4]
			} else {
				repos = append(repos[:i], repos[i+1:]...)

			}
			break
		}
	}
	repos = append(r, repos...)
	r = r[:5]
	repoJSON, err := json.MarshalIndent(repos, "", "  ")
	check(err)
	json, err := os.Open(filePath)
	json.Write(repoJSON)
	check(err)
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
