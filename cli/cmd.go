package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/khaledAlorayir/yts-cli/yts"
)

type searchMoviesMsg struct {
	movies []yts.Option
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func searchMovies(query string) tea.Cmd {
	return func() tea.Msg {
		movies, err := yts.GetMovies(query)

		if err != nil {
			return errMsg{err: err}
		}

		return searchMoviesMsg{movies: movies}
	}
}
