package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/khaledAlorayir/yts-cli/yts"
)

type step int

const (
	SEARCH_INPUT step = iota
	MOVIE_LIST
	VERSION_LIST
)

type model struct {
	textInput     textinput.Model
	movies        []yts.Option
	movieVersions []yts.Option
	selectedIndex int
	err           error
	step          step
}

func InitialModel() model {
	textInput := textinput.New()
	textInput.Placeholder = "movie name"
	textInput.Focus()
	textInput.CharLimit = 156

	return model{textInput: textInput, step: SEARCH_INPUT}
}

func (model model) Init() tea.Cmd {
	return textinput.Blink
}

func (model model) View() string {
	if model.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", model.err)
	}

	if model.step == SEARCH_INPUT {
		return fmt.Sprintf(
			"Please enter the movie name\n\n%s\n\n%s\n",
			model.textInput.View(),
			"(esc to quit)",
		) + "\n"
	}

	if model.step == MOVIE_LIST || model.step == VERSION_LIST {
		var labels []string
		var options []yts.Option

		if model.step == MOVIE_LIST {
			options = model.movies
		} else {
			options = model.movieVersions
		}

		for i, o := range options {
			if i == model.selectedIndex {
				labels = append(labels, fmt.Sprintf("-> %s", o.Label))
			} else {
				labels = append(labels, fmt.Sprintf("   %s", o.Label))
			}
		}

		return fmt.Sprintf("   Pick a movie!\n\n%s", strings.Join(labels, "\n"))
	}
	return ""
}

func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return model, tea.Quit
		case tea.KeyEnter:
			return model, model.handleStep()
		case tea.KeyUp, tea.KeyDown:
			return model.handleArrows(msg), nil
		}
	case searchMoviesMsg:
		model.step = MOVIE_LIST
		model.movies = msg.movies
		return model, nil

	case searchVersionsMsg:
		model.step = VERSION_LIST
		model.selectedIndex = 0
		model.movieVersions = msg.versions
		return model, nil

	case errMsg:
		model.err = msg
		return model, tea.Quit
	}

	model.textInput, cmd = model.textInput.Update(msg)
	return model, cmd
}

func (model model) handleArrows(msg tea.KeyMsg) model {
	switch msg.Type {
	case tea.KeyUp:
		model.selectedIndex--
	case tea.KeyDown:
		model.selectedIndex++
	}

	optionsCount := 0

	if model.step == MOVIE_LIST {
		optionsCount = len(model.movies)
	} else {
		optionsCount = len(model.movieVersions)
	}

	model.selectedIndex = (model.selectedIndex + optionsCount) % optionsCount
	return model
}

func (model model) handleStep() tea.Cmd {
	switch model.step {
	case SEARCH_INPUT:
		return searchMovies(model.textInput.Value())
	case MOVIE_LIST:
		return searchVersions(model.movies[model.selectedIndex].Url)
	default:
		return nil
	}
}
