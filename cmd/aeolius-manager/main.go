package main

import (
	"flag"
	"log"

	"github.com/pojntfx/aeolius/pkg/persisters"
)

func main() {
	postgresUrl := flag.String("postgres-url", "postgresql://postgres@localhost:5432/aeolius?sslmode=disable", "PostgreSQL URL")

	flag.Parse()

	persister := persisters.NewManagerPersister(*postgresUrl)

	if err := persister.Open(); err != nil {
		panic(err)
	}
	defer persister.Close()

	log.Println("Connected to PostgreSQL")
}
