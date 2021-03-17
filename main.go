package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/NinjaAung/traverse/scraping"
)

func main() {

	filePath := "/tmp/traverse_recent.json"

	if runtime.GOOS == "windows" {
		user, _ := user.Current()
		filePath = filepath.Join(user.HomeDir, "appData/Local/Temp/traverse_recent.json")
	}

	if len(os.Args) == 2 {
		repo, err := scraping.Run(os.Args[1])
		if err != nil {
			panic(err)
		}
		fmt.Println(repo)
		repo.SaveToJSON(filePath)
		os.Exit(0)
	}
}
