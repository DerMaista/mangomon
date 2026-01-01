package tools

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TransformSelectedMsg struct {
	Transform int
}

type TransformCancelledMsg struct{}

type TransformPickerModel struct {
	Monitor   string
	Selected  int
	CurrentID int
}

var transformNames = []string{
	"Normal (0)",
	"90° (1)",
	"180° (2)",
	"270° (3)",
	"Flipped (4)",
	"Flipped 90° (5)",
	"Flipped 180° (6)",
	"Flipped 270° (7)",
}

func NewTransformPicker(monitor string, currentVal int) TransformPickerModel {
	return TransformPickerModel{
		Monitor:   monitor,
		Selected:  currentVal,
		CurrentID: currentVal,
	}
}

func (m TransformPickerModel) Init() tea.Cmd {
	return nil
}

func (m TransformPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, func() tea.Msg { return TransformCancelledMsg{} }
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			}
		case "down", "j":
			if m.Selected < len(transformNames)-1 {
				m.Selected++
			}
		case "enter":
			return m, func() tea.Msg { return TransformSelectedMsg{Transform: m.Selected} }
		}
	}
	return m, nil
}

func (m TransformPickerModel) View() string {
	s := fmt.Sprintf("Select Rotation for %s\n\n", m.Monitor)

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	normalStyle := lipgloss.NewStyle().PaddingLeft(2)

	for i, name := range transformNames {
		cursor := "  "
		if i == m.Selected {
			cursor = "▶ "
		}

		line := name
		if i == m.CurrentID {
			line += " (current)"
		}

		if i == m.Selected {
			s += selectedStyle.Render(cursor+line) + "\n"
		} else {
			s += normalStyle.Render(line) + "\n"
		}
	}

	s += "\n[Enter] Select  [Esc] Cancel"

	return s
}
