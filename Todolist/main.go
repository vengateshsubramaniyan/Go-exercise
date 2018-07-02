package main

import (
	"Go-exercise/Todolist/cmd"
	"Go-exercise/Todolist/db"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	dirName, err := homedir.Dir()
	taskbucket.ErrorReporter(err)
	taskbucket.Init(filepath.Join(dirName, "my.db"))
	cmd.Execute()
}
