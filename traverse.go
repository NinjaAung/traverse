package traverse

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

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

func isFileExists(filePath string) error {
	_, err := os.Open(filePath)
	return err
}

// SaveRepo saves the Repo object to a json file
func SaveRepo(filePath string, repo *Repo) {
	if isFileExists(filePath) != nil {
		f, err := os.Create(filePath)
		check(err)
		repoJSON, err := json.MarshalIndent([]*Repo{repo}, "", "  ")
		check(err)
		f.Write(repoJSON)
		f.Close()
	} else {
		updateJSON(filePath, repo)
	}
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
		repos = repos[:5-fillerSize]
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

// InitRepo ...
func InitRepo(baseURL string, repo *Repo) {
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

// SearchFolder ...
func SearchFolder(link string, dir *Dir) {
	c := colly.NewCollector(colly.Async(true))
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})
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
			SearchFolder("https://github.com"+href, &newDir)
			dir.Dirs = append(dir.Dirs, newDir)
		}

	})
	c.Visit(link)
	c.Wait()
}
