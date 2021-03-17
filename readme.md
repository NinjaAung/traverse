<p align="center">
Traverse
<br>
<br>
Traverse any git repo with go!
</p>
<p align="center">
  <a>
    <a href="https://goreportcard.com/badge/github.com/NinjaAung/traverse" />
    <img alt="commits" src="https://goreportcard.com/badge/github.com/NinjaAung/traverse" target="_blank" />
    <a href="https://github.com/NinjaAung/NinjaAung/commits/master">
    <img alt="commits" src="https://img.shields.io/github/commit-activity/w/NinjaAung/traverse?color=green" target="_blank" />
  </a> 
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>
<br>

## 🚀  Getting Started\

### Description

In short this packages make it so that users can grab teh contents of github repo and 

### Installation

Package can be downloaded by using go

```bash
go get github.com/NinjaAung/traverse
```

### Usage

There are three main function that traverse has:

#### Scraping.Run(repoName)

Scraping.Run() is of course resposible of returning a repo struct resposible with data: Name, Route, File and Dir used like this:

```go
repo, _ := Scraping.Run("NinjaAung/Traverse")
```

#### *repo.SaveToJSON(filePath)

SaveToJSON takes a repo struct and marshall it's data in a list of repo in json, used like this:

```go
repo.SaveToJSON("example.json")
```

#### tra(filePath)

tra, short for traverse take filePath to the json info and reads the first item or any if desire and allows the user to traverse through it in the terminal

```go
tra("example.json")
```


## Misc

Here is a collection of traverse related items:

- General landscaping / future plans with [Traverse](https://www.figma.com/file/4IgIZ1sVTaL1eCpITjHmIt?embed_host=share&kind=&node-id=0%3A1&viewer=1)
- [Google Slides and Recording](https://docs.google.com/presentation/d/1kS4SOqg5-zKV0U8FgTJvd3fbHrSrNxU2hOxzWDTo6rA/edit?usp=sharing)