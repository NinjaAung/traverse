package main

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/NinjaAung/traverse/scraping"
	"github.com/NinjaAung/traverse/traverse"
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
		repo.SaveToJSON(filePath)
		traverse.Tra(repo.Dir)
		os.Exit(0)
	}
	traverse.ReadRecent(filePath)
	
	os.Exit(0)
}
