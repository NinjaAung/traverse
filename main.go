package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/NinjaAung/traverse/traverse"
)

func main() {
	// repo, err := scraping.Run("Make-School-Courses/SPD-2.31-Testing-and-Architecture")
	// if err != nil {
	// 	panic(err)
	// }
	// repo.SaveToJSON("test.json")
	f, _ := ioutil.ReadFile("test.json")
	var repos []traverse.Repo
	json.Unmarshal(f, &repos)

	fmt.Println(repos[0].Dir.Name)

	tra(repos[0].Dir)

}

func tra(dir traverse.Dir) func() {
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
			return tra(dir)
		}
		fmt.Println("Went back")

	} else if optionNum < int64(fileStart)+1 {
		fmt.Printf("Folder: %s\n", dir.Dirs[int(optionNum)-1].Name)
		return tra(*dir.Dirs[int(optionNum)-1])

	} else if optionNum >= int64(fileStart)+1 {
		index := int(optionNum) - 1
		fmt.Printf("File: %s", dir.Files[index-fileStart])
		return tra(dir)

	}
	exec.Command("clear").Run()
	return tra(dir)
}

func raw(path string) string {
	return "https://raw.githubusercontent.com/" + path
}
