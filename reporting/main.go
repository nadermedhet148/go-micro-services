package main

import (
	"example.com/reporting/consumers"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	consumers.RunGroup()
}
