package cli

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khaledAlorayir/yts-cli/common"
	"github.com/khaledAlorayir/yts-cli/services/downloader"
	"github.com/khaledAlorayir/yts-cli/services/yts"
)

type searchMoviesMsg struct {
	movies []common.Option
}

type searchVersionsMsg struct {
	versions []common.Option
}

type goToStepMsg struct {
	step step
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func searchMovies(query string) tea.Cmd {
	return func() tea.Msg {
		movies, err := yts.GetMovies(query)

		if err != nil {
			return errMsg{err: err}
		}

		if movies == nil {
			return errMsg{err: errors.New("no movies found ;(")}
		}

		return searchMoviesMsg{movies: movies}
	}
}

func searchVersions(link string) tea.Cmd {
	return func() tea.Msg {
		versions, err := yts.GetMovieVersionOptions(link)

		if err != nil {
			return errMsg{err: err}
		}

		return searchVersionsMsg{versions: versions}
	}
}

func goToStep(step step) tea.Cmd {
	return func() tea.Msg {
		return goToStepMsg{step: step}
	}
}

func downloadMovie(movie common.Option, movieTitle string) tea.Cmd {
	return func() tea.Msg {
		err := downloader.SaveFile(movie.Url, fmt.Sprintf("%s_%s", movieTitle, movie.Label))

		if err != nil {
			return errMsg{err: err}
		}

		return goToStepMsg{step: MOVIE_DOWNLOADED}
	}
}
