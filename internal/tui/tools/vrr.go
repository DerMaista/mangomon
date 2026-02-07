package tools

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type VRRSelectedMsg struct {
	VRR int
}

type VRRCancelledMsg struct{}

type VRRPickerModel struct {
	Monitor   string
	Selected  int
	CurrentID int
}

var vrrNames = []string{
	"Disabled (0)",
	"Enabled (1)",
}

func NewVRRPicker(monitor string, currentVal int) VRRPickerModel {
	selected := currentVal
	if selected < 0 || selected >= len(vrrNames) {
		selected = 0
	}
	return VRRPickerModel{
		Monitor:   monitor,
		Selected:  selected,
		CurrentID: currentVal,
	}
}

func (m VRRPickerModel) Init() tea.Cmd {
	return nil
}

func (m VRRPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, func() tea.Msg { return VRRCancelledMsg{} }
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			}
		case "down", "j":
			if m.Selected < len(vrrNames)-1 {
				m.Selected++
			}
		case "enter":
			return m, func() tea.Msg { return VRRSelectedMsg{VRR: m.Selected} }
		}
	}
	return m, nil
}

func (m VRRPickerModel) View() string {
	s := fmt.Sprintf("Variable Refresh Rate for %s\n\n", m.Monitor)

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
	normalStyle := lipgloss.NewStyle().PaddingLeft(2)

	for i, name := range vrrNames {
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
