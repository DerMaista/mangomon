package tools

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MirrorSelectedMsg struct {
	TargetID string
}

type MirrorCancelledMsg struct{}

type MirrorPickerModel struct {
	SourceID string
	Targets  []string
	Selected int
}

func NewMirrorPicker(source string, allMonitors []string) MirrorPickerModel {
	var targets []string
	// Filter out self
	for _, m := range allMonitors {
		if m != source {
			targets = append(targets, m)
		}
	}

	return MirrorPickerModel{
		SourceID: source,
		Targets:  targets,
		Selected: 0,
	}
}

func (m MirrorPickerModel) Init() tea.Cmd {
	return nil
}

func (m MirrorPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, func() tea.Msg { return MirrorCancelledMsg{} }
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			}
		case "down", "j":
			if m.Selected < len(m.Targets)-1 {
				m.Selected++
			}
		case "enter":
			if len(m.Targets) > 0 {
				return m, func() tea.Msg { return MirrorSelectedMsg{TargetID: m.Targets[m.Selected]} }
			}
		}
	}
	return m, nil
}

func (m MirrorPickerModel) View() string {
	s := fmt.Sprintf("Select Monitor to Mirror %s to:\n\n", m.SourceID)

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	normalStyle := lipgloss.NewStyle().PaddingLeft(2)

	if len(m.Targets) == 0 {
		return s + "No other monitors available to mirror."
	}

	for i, target := range m.Targets {
		cursor := "  "
		if i == m.Selected {
			cursor = "▶ "
		}

		line := target

		if i == m.Selected {
			s += selectedStyle.Render(cursor+line) + "\n"
		} else {
			s += normalStyle.Render(line) + "\n"
		}
	}

	s += "\n[Enter] Mirror  [Esc] Cancel"

	return s
}
