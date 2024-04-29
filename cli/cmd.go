package cli

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khaledAlorayir/yts-cli/yts"
)

type searchMoviesMsg struct {
	movies []yts.Option
}

type searchVersionsMsg struct {
	versions []yts.Option
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
