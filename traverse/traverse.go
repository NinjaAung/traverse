package traverse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Dir is an embedded struct of folders in a repo
type Dir struct {
	Name  string
	Route string
	Files []string
	Dirs  []*Dir
}

//NewDir ...
func NewDir(name, route string) Dir {
	return Dir{
		Name:  name,
		Route: route,
		Files: []string{},
		Dirs:  []*Dir{},
	}
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

// SaveToJSON saves the Repo object to a json file
func (repo *Repo) SaveToJSON(filePath string) {
	if isFileExists(filePath) != nil {
		fmt.Println("Creating File")
		f, err := os.Create(filePath)
		check(err)
		repoJSON, err := json.MarshalIndent([]*Repo{repo}, "", "  ")
		check(err)
		f.Write(repoJSON)
		f.Close()
	} else {
		fmt.Println("Updating File")
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
