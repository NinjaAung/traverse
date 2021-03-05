package main

import (
	"github.com/NinjaAung/traverse/scraping"
)

func main() {
	repo, err := scraping.Run("Make-School-Courses/SPD-2.31-Testing-and-Architecture")
	if err != nil {
		panic(err)
	}
	repo.SaveToJSON("test.json")
}
