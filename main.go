package main

import (
	"fmt"

	"github.com/khaledAlorayir/yts-cli/services"
)

func main() {
	fmt.Println(services.GetMovieOptions("https://yts.mx/movies/raiders-of-the-lost-ark-1981"))
	// fmt.Println(services.GetMovies("raiders of the lost ark"))
}
