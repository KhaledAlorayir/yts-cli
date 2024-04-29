package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/khaledAlorayir/yts-cli/yts"
)

type model struct {
	textInput     textinput.Model
	movies        []yts.Option
	selectedIndex int
	err           error
}

func InitialModel() model {
	textInput := textinput.New()
	textInput.Placeholder = "movie name"
	textInput.Focus()
	textInput.CharLimit = 156
	return model{textInput: textInput}
}

func (model model) Init() tea.Cmd {
	return textinput.Blink
}

func (model model) View() string {
	if model.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", model.err)
	}

	if len(model.movies) > 0 {
		var options []string

		for i, o := range model.movies {
			if i == model.selectedIndex {
				options = append(options, fmt.Sprintf("-> %s", o.Label))
			} else {
				options = append(options, fmt.Sprintf("   %s", o.Label))
			}
		}

		return strings.Join(options, "\n")
	}

	return fmt.Sprintf(
		"Please enter the movie name\n\n%s\n\n%s\n",
		model.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return model, tea.Quit
		case tea.KeyEnter:
			return model, searchMovies(model.textInput.Value())
		case tea.KeyUp, tea.KeyDown:
			return model.handleArrows(msg), nil
		}
	case searchMoviesMsg:
		model.movies = msg.movies
		return model, nil

	case errMsg:
		model.err = msg
		return model, tea.Quit
	}
	// We handle errors just like any other message

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

	moviesCount := len(model.movies)
	model.selectedIndex = (model.selectedIndex + moviesCount) % moviesCount
	return model
}
