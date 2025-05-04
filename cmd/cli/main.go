package main

import (
	"github.com/up1io/muxo/cli"
	"log"
)

func main() {
	r := cli.New()
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run cli. %s", err)
	}
}
