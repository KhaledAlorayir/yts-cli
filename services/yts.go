package services

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

const (
	base_url = "https://yts.mx/browse-movies/%s/all/all/0/latest/0/all"
)

func GetMovies(query string) ([]Option, error) {
	c := colly.NewCollector()
	url := fmt.Sprintf(base_url, query)

	var movies []Option
	var err error

	c.OnError(func(_ *colly.Response, CollectorErr error) {
		log.Println("Something went wrong: ", CollectorErr)
		err = errors.New("something went wrong")
	})

	c.OnHTML(".browse-movie-bottom", func(e *colly.HTMLElement) {
		title := e.ChildText(".browse-movie-title")
		year := e.ChildText(".browse-movie-year")
		url := e.ChildAttr(".browse-movie-title", "href")

		movies = append(movies, Option{label: fmt.Sprintf("%s - %s", title, year), url: url})
	})

	c.Visit(url)

	if err != nil {
		return movies, err
	}

	return movies, nil
}

func GetMovieOptions(link string) ([]Option, error) {
	c := colly.NewCollector()

	var links []Option
	var err error

	c.OnError(func(_ *colly.Response, CollectorErr error) {
		log.Println("something went wrong: ", err)
		err = errors.New("something went wrong")
	})

	c.OnHTML("#movie-info .hidden-sm a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if strings.HasPrefix(url, "https://yts.mx/torrent") {
			links = append(links, Option{label: e.Text, url: url})
		}
	})
	c.Visit(link)

	if err != nil {
		return links, err
	}

	return links, nil
}
