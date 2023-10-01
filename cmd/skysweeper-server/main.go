package main

import "github.com/pojntfx/skysweeper/cmd/skysweeper-server/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
