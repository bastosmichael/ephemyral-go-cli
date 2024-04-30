package main

import (
	"ephemyral/cmd"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "gendoc" {
		cmd.GenerateMarkdownDocs()
		return
	}
	cmd.Execute()
}
