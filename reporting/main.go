package main

import (
	"example.com/reporting/consumers"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	go consumers.RunGroup()
	go consumers.RunStream()
	select {} // wait forever
}
