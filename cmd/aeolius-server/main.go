package main

import "github.com/pojntfx/aeolius/cmd/aeolius-server/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
