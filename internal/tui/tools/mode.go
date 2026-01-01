package tools

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Mode struct {
	Width, Height int
	Rate          float64
}

func (m Mode) String() string {
	return fmt.Sprintf("%dx%d @ %.2fHz", m.Width, m.Height, m.Rate)
}

type ModeSelectedMsg struct {
	Mode Mode
}

type ModeCancelledMsg struct{}

type ModePickerModel struct {
	Monitor  string
	Modes    []Mode
	Selected int
	Current  Mode
}

func NewModePicker(monitor string, current Mode, modes []Mode) ModePickerModel {
	return ModePickerModel{
		Monitor:  monitor,
		Modes:    modes,
		Selected: 0,
		Current:  current,
	}
}

func (m ModePickerModel) Init() tea.Cmd {
	return nil
}

func (m ModePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, func() tea.Msg { return ModeCancelledMsg{} }
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			}
		case "down", "j":
			if m.Selected < len(m.Modes)-1 {
				m.Selected++
			}
		case "home", "g":
			m.Selected = 0
		case "end", "G":
			m.Selected = len(m.Modes) - 1
		case "enter":
			return m, func() tea.Msg { return ModeSelectedMsg{Mode: m.Modes[m.Selected]} }
		}
	}
	return m, nil
}

func (m ModePickerModel) View() string {
	s := fmt.Sprintf("Select Mode for %s\n\n", m.Monitor)

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	normalStyle := lipgloss.NewStyle().PaddingLeft(2)
	currentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))

	for i, mode := range m.Modes {
		cursor := "  "
		if i == m.Selected {
			cursor = "▶ "
		}

		modeStr := mode.String()
		indicator := ""
		if mode.Width == m.Current.Width && mode.Height == m.Current.Height && mode.Rate == m.Current.Rate {
			indicator = currentStyle.Render(" (current)")
		}

		line := fmt.Sprintf("%s%s", modeStr, indicator)

		if i == m.Selected {
			s += selectedStyle.Render(cursor+line) + "\n"
		} else {
			s += normalStyle.Render(line) + "\n"
		}
	}

	s += "\n[Enter] Select  [Esc] Cancel"

	return s
}
