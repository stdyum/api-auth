package main

import (
	"log"

	"github.com/stdyum/api-auth/internal"
)

func main() {
	log.Fatalf("error launching web server %s", internal.App().Error())
}
