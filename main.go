package main

import (
	"fmt"

	"github.com/khaledAlorayir/yts-cli/services"
)

func main() {
	thing, _ := services.GetMovies("007")
	fmt.Println(len(thing))
}
