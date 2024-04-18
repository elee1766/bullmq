package src

import (
	"embed"
	"io/fs"
)

//go:embed commands
var commands embed.FS

var Commands fs.FS

func init() {
	loaded, err := fs.Sub(commands, "commands")
	if err != nil {
		panic(err)
	}
	Commands = loaded
}
