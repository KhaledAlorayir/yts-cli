package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	SEARCH_INPUT step = iota
	MOVIE_LIST
	VERSION_LIST
	MOVIE_DOWNLOADED
)

type model struct {
	textInput     textinput.Model
	movies        cliList
	movieVersions cliList
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

	ui := ""

	switch model.step {
	case SEARCH_INPUT:
		ui += fmt.Sprintf("Please enter the movie name\n\n%s", model.textInput.View())
	case MOVIE_LIST:
		ui += model.movies.view()
	case VERSION_LIST:
		ui += model.movieVersions.view()
	case MOVIE_DOWNLOADED:
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
			if model.step == MOVIE_LIST || model.step == VERSION_LIST {
				return model.handleArrows(msg), nil
			}
			return model, nil
		}
	case searchMoviesMsg:
		model.movies = newCliList(msg.movies)
		return model, goToStep(MOVIE_LIST)

	case searchVersionsMsg:
		model.movieVersions = newCliList(msg.versions)
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
	switch model.step {
	case MOVIE_LIST:
		model.movies.setSelection(msg)
	case VERSION_LIST:
		model.movieVersions.setSelection(msg)
	}
	return model
}

func (model model) handleSelection() tea.Cmd {
	switch model.step {
	case SEARCH_INPUT:
		//TODO input validation here
		return searchMovies(model.textInput.Value())
	case MOVIE_LIST:
		if model.movies.isPreviousOptionSelected() {
			return goToStep(SEARCH_INPUT)
		}
		return searchVersions(model.movies.getSelected().Url)
	case VERSION_LIST:
		if model.movieVersions.isPreviousOptionSelected() {
			return goToStep(MOVIE_LIST)
		}
		return downloadMovie(model.movieVersions.getSelected(), model.movies.getSelected().Label)
	default:
		return nil
	}
}
