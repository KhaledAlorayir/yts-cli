package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/khaledAlorayir/yts-cli/common"
)

type step int

const (
	SEARCH_INPUT step = iota
	MOVIE_LIST
	VERSION_LIST
	MOVIE_DOWNLOADED
)

type model struct {
	textInput             textinput.Model
	movies                []common.Option
	movieVersions         []common.Option
	selectedMovieIndex    int
	selectedVersionsIndex int
	err                   error
	step                  step
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

	ui := ""

	if model.step == SEARCH_INPUT {
		ui += fmt.Sprintf(
			"Please enter the movie name\n\n%s",
			model.textInput.View(),
		)
	}

	if model.step == MOVIE_LIST || model.step == VERSION_LIST {
		var labels []string
		var options []common.Option
		var selectedIndex int

		if model.step == MOVIE_LIST {
			options = model.movies
			selectedIndex = model.selectedMovieIndex

		} else {
			options = model.movieVersions
			selectedIndex = model.selectedVersionsIndex
		}

		for i, o := range options {
			var label string

			if i == selectedIndex {
				label = fmt.Sprintf("-> %s", o.Label)
			} else {
				label = fmt.Sprintf("   %s", o.Label)
			}

			if i == len(options)-1 {
				label = "\n" + label
			}

			labels = append(labels, label)
		}

		ui += fmt.Sprintf("   Pick an option!\n\n%s", strings.Join(labels, "\n"))
	}

	if model.step == MOVIE_DOWNLOADED {
		ui += "torrent has been saved to your downloads folder!"
	}

	return ui + fmt.Sprintf("\n\n%s\n\n", "(esc to quit)")
}

func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return model, tea.Quit
		case tea.KeyEnter:
			return model, model.handleSelection()
		case tea.KeyUp, tea.KeyDown:
			return model.handleArrows(msg), nil
		}
	case searchMoviesMsg:
		model.movies = msg.movies
		return model, goToStep(MOVIE_LIST)

	case searchVersionsMsg:
		model.movieVersions = msg.versions
		return model, goToStep(VERSION_LIST)

	case goToStepMsg:
		model.step = msg.step
		return model, nil

	case errMsg:
		model.err = msg
		return model, tea.Quit
	}

	model.textInput, cmd = model.textInput.Update(msg)
	return model, cmd
}

func (model model) handleArrows(msg tea.KeyMsg) model {
	var selectedIndex int
	var optionsCount int

	if model.step == MOVIE_LIST {
		selectedIndex = model.selectedMovieIndex
		optionsCount = len(model.movies)

	} else {
		selectedIndex = model.selectedVersionsIndex
		optionsCount = len(model.movieVersions)
	}

	switch msg.Type {
	case tea.KeyUp:
		selectedIndex--
	case tea.KeyDown:
		selectedIndex++
	}

	if model.step == MOVIE_LIST {
		model.selectedMovieIndex = (selectedIndex + optionsCount) % optionsCount
	} else {
		model.selectedVersionsIndex = (selectedIndex + optionsCount) % optionsCount
	}

	return model
}

func (model model) handleSelection() tea.Cmd {
	switch model.step {
	case SEARCH_INPUT:
		//TODO input validation here
		return searchMovies(model.textInput.Value())
	case MOVIE_LIST:
		if model.selectedMovieIndex == len(model.movies)-1 {
			return goToStep(SEARCH_INPUT)
		}
		return searchVersions(model.movies[model.selectedMovieIndex].Url)
	case VERSION_LIST:
		if model.selectedVersionsIndex == len(model.movieVersions)-1 {
			return goToStep(MOVIE_LIST)
		}
		return downloadMovie(model.movieVersions[model.selectedVersionsIndex], model.movies[model.selectedMovieIndex].Label)
	default:
		return nil
	}
}
