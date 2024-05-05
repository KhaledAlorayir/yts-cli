package yts

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/khaledAlorayir/yts-cli/common"
)

const (
	base_url  = "https://yts.mx/browse-movies/%s/all/all/0/latest/0/all"
	max_pages = 1
)

func GetMovies(query string) ([]common.Option, error) {
	url := fmt.Sprintf(base_url, query)

	c := colly.NewCollector()

	var movies []common.Option
	var err error
	page := 1

	c.OnError(func(_ *colly.Response, CollectorErr error) {
		err = errors.New("something went wrong")
	})

	c.OnHTML(".browse-movie-bottom", func(e *colly.HTMLElement) {
		title := e.ChildText(".browse-movie-title")
		year := e.ChildText(".browse-movie-year")
		url := e.ChildAttr(".browse-movie-title", "href")

		movies = append(movies, common.Option{Label: fmt.Sprintf("%s - %s", title, year), Url: url})
	})

	c.OnHTML("section + .hidden-sm .tsc_pagination li a", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Text, "Next") && page < max_pages {
			page++
			e.Request.Visit(e.Attr("href"))
		}
	})

	err = c.Visit(url)

	if err != nil {
		return movies, err
	}

	return movies, nil
}

func GetMovieVersionOptions(link string) ([]common.Option, error) {
	c := colly.NewCollector()

	var links []common.Option
	var err error

	c.OnError(func(_ *colly.Response, CollectorErr error) {
		err = errors.New("something went wrong")
	})

	c.OnHTML("#movie-info .hidden-sm a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if strings.HasPrefix(url, "https://yts.mx/torrent") {
			links = append(links, common.Option{Label: e.Text, Url: url})
		}
	})

	err = c.Visit(link)

	if err != nil {
		return links, err
	}

	return links, nil
}
