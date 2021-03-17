package traverse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
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

func clear() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()

}

//ReadRecent
func ReadRecent() {

}

//Tra ...
func Tra(dir Dir) func() {
	reader := bufio.NewReader(os.Stdin)
	fileStart := 0

	fmt.Printf("0. Back\nFolders:\n\n")
	for i, dir := range dir.Dirs {
		fmt.Printf("%d. %s\n", i+1, dir.Name)
		fileStart = i + 1
	}
	fmt.Printf("\nFiles:\n\n")
	for i, file := range dir.Files {
		fmt.Printf("%d. %s\n", fileStart+i+1, file)
	}
	fmt.Print(": ")
	option, _ := reader.ReadString('\n')
	optionNum, _ := strconv.ParseInt(strings.TrimSpace(option), 10, 64)
	if optionNum == 0 {
		if dir.Route == "master" {
			fmt.Print("This is root")
			clear()
			return Tra(dir)
		}
		fmt.Println("Went back")

	} else if optionNum < int64(fileStart)+1 {
		fmt.Printf("Folder: %s\n", dir.Dirs[int(optionNum)-1].Name)
		clear()
		return Tra(*dir.Dirs[int(optionNum)-1])

	} else if optionNum >= int64(fileStart)+1 {
		index := int(optionNum) - 1
		fmt.Printf("File: %s", dir.Files[index-fileStart])
		clear()
		return Tra(dir)

	}
	clear()
	return Tra(dir)
}
