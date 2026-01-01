package tools

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LayoutSelectedMsg struct {
	Layout string
}

type LayoutCancelledMsg struct{}

type LayoutPickerModel struct {
	Monitor       string
	Selected      int
	CurrentLayout string
}

var layoutOptions = []string{
	"tile",
	"scroller",
	"grid",
	"deck",
	"monocle",
	"center_tile",
	"vertical_tile",
	"vertical_scroller",
}

func NewLayoutPicker(monitor string, currentLayout string) LayoutPickerModel {
	selectedIdx := 0
	for i, l := range layoutOptions {
		if l == currentLayout {
			selectedIdx = i
			break
		}
	}

	return LayoutPickerModel{
		Monitor:       monitor,
		Selected:      selectedIdx,
		CurrentLayout: currentLayout,
	}
}

func (m LayoutPickerModel) Init() tea.Cmd {
	return nil
}

func (m LayoutPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, func() tea.Msg { return LayoutCancelledMsg{} }
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			}
		case "down", "j":
			if m.Selected < len(layoutOptions)-1 {
				m.Selected++
			}
		case "enter":
			return m, func() tea.Msg { return LayoutSelectedMsg{Layout: layoutOptions[m.Selected]} }
		}
	}
	return m, nil
}

func (m LayoutPickerModel) View() string {
	s := fmt.Sprintf("Select Layout for %s\n\n", m.Monitor)

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	normalStyle := lipgloss.NewStyle().PaddingLeft(2)

	for i, name := range layoutOptions {
		cursor := "  "
		if i == m.Selected {
			cursor = "▶ "
		}

		line := name
		if name == m.CurrentLayout {
			line += " (current)"
		}

		if i == m.Selected {
			s += selectedStyle.Render(cursor+line) + "\n"
		} else {
			s += normalStyle.Render(line) + "\n"
		}
	}

	s += "\n[Enter] Select  [Esc] Cancel"
	s += "\n\n⚠ Note: Restart MangoWC to apply global monitor layouts"

	return s
}
