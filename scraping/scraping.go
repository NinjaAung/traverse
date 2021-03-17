package scraping

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/NinjaAung/traverse/traverse"
	"github.com/gocolly/colly"
)

var (
	foldersChan = make(chan *traverse.Dir)
	wg          sync.WaitGroup
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Run ...
func Run(repoName string) (traverse.Repo, error) {
	resp, _ := http.Get("https://github.com/" + repoName)
	if resp.StatusCode == 404 {
		return traverse.Repo{}, fmt.Errorf("repo dosen't exsist, spelled incorrectly, or repo is private")
	}
	for i := 0; i < 100; i++ {
		go searchworker(foldersChan)
	}
	repo := traverse.Repo{Name: repoName}
	initRepo(repoName, &repo)
	dir := traverse.NewDir(repoName, repoName)
	start := time.Now()
	wg.Add(1)
	foldersChan <- &dir
	wg.Wait()
	close(foldersChan)
	repo.Dir = dir
	fmt.Println("Repo Collected took: ", time.Since(start))
	return repo, nil
}

func searchworker(foldersChan chan *traverse.Dir) {
	for folder := range foldersChan {
		directory := SearchFolder("https://github.com/"+folder.Route, folder)
		for i := range directory.Dirs {
			wg.Add(1)
			foldersChan <- directory.Dirs[i]
		}
		wg.Done()
	}
}

//SearchFolder ...
func SearchFolder(link string, dir *traverse.Dir) *traverse.Dir {
	c := colly.NewCollector()
	c.OnHTML("span a.js-navigation-open", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		link := strings.Split(href, "/")
		dir.Route = strings.Join(link[4:len(link)-1], "/")
		routeType := link[3]
		title := e.Attr("title")

		if routeType == "tree" {
			newDir := traverse.NewDir(title, href)
			dir.Dirs = append(dir.Dirs, &newDir)
		} else if routeType == "blob" {
			dir.Files = append(dir.Files, title)
		}

	})
	c.Visit(link)
	return dir
}

func initRepo(baseURL string, repo *traverse.Repo) {
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
	c.Visit("https://github.com/" + baseURL)
}
