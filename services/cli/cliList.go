package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khaledAlorayir/yts-cli/common"
)

type cliList struct {
	options       []common.Option
	selectedIndex int
}

func newCliList(options []common.Option) cliList {
	options = append(options, common.Option{Label: "previous"})

	return cliList{options: options}
}

func (list *cliList) setSelection(msg tea.KeyMsg) {
	optionsCount := len(list.options)

	switch msg.Type {
	case tea.KeyUp:
		list.selectedIndex--
	case tea.KeyDown:
		list.selectedIndex++
	}

	list.selectedIndex = (list.selectedIndex + optionsCount) % optionsCount
}

func (list cliList) view() string {
	var labels []string

	for i, o := range list.options {
		var label string

		if i == list.selectedIndex {
			label = fmt.Sprintf("-> %s", o.Label)
		} else {
			label = fmt.Sprintf("   %s", o.Label)
		}

		if i == len(list.options)-1 {
			label = "\n" + label
		}

		labels = append(labels, label)
	}

	return fmt.Sprintf("   Pick an option!\n\n%s", strings.Join(labels, "\n"))
}

func (list cliList) isPreviousOptionSelected() bool {
	return list.selectedIndex == len(list.options)-1
}

func (list cliList) getSelected() common.Option {
	return list.options[list.selectedIndex]
}
