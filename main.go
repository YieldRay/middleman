package main

import (
	"fmt"
	"os"

	"github.com/yieldray/middleman/cmd"
	"github.com/yieldray/middleman/gui"
)

func main() {
	if len(os.Args) == 1 {
		cwd, _ := os.Getwd()
		fmt.Printf("cwd: %s\n", cwd)
		gui.Main()
	} else {
		cmd.Execute()
	}
}
