package tools

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ScaleSelectedMsg struct {
	Scale float64
}

type ScaleCancelledMsg struct{}

type ScalePickerModel struct {
	Monitor  string
	Scales   []float64
	Selected int
	Current  float64

	CustomMode  bool
	CustomInput textinput.Model

	Width, Height int
}

func NewScalePicker(monitor string, current float64, width, height int) ScalePickerModel {
	ti := textinput.New()
	ti.Placeholder = "1.0"
	ti.CharLimit = 5
	ti.Width = 10

	return ScalePickerModel{
		Monitor:     monitor,
		Scales:      []float64{0.5, 0.75, 1.0, 1.25, 1.5, 1.75, 2.0, 2.5, 3.0},
		Selected:    2,
		Current:     current,
		CustomInput: ti,
		Width:       width,
		Height:      height,
	}
}

func (m ScalePickerModel) Init() tea.Cmd {
	return nil
}

func (m ScalePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Custom Input Mode
	if m.CustomMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.CustomMode = false
				m.CustomInput.Blur()
				return m, nil
			case "enter":
				val := m.CustomInput.Value()
				if s, err := strconv.ParseFloat(val, 64); err == nil && s > 0 && s <= 10 {
					return m, func() tea.Msg { return ScaleSelectedMsg{Scale: s} }
				}
			}
		}
		m.CustomInput, cmd = m.CustomInput.Update(msg)
		return m, cmd
	}

	// List Mode
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, func() tea.Msg { return ScaleCancelledMsg{} }
		case "c":
			m.CustomMode = true
			m.CustomInput.Focus()
			return m, textinput.Blink
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			}
		case "down", "j":
			if m.Selected < len(m.Scales)-1 {
				m.Selected++
			}
		case "home", "g":
			m.Selected = 0
		case "end", "G":
			m.Selected = len(m.Scales) - 1
		case "enter":
			return m, func() tea.Msg { return ScaleSelectedMsg{Scale: m.Scales[m.Selected]} }

		// Quick keys
		case "1":
			return m, func() tea.Msg { return ScaleSelectedMsg{Scale: 1.0} }
		case "2":
			return m, func() tea.Msg { return ScaleSelectedMsg{Scale: 2.0} }
		}
	}

	return m, nil
}

func (m ScalePickerModel) View() string {
	if m.CustomMode {
		return fmt.Sprintf(
			"Enter Custom Scale for %s:\n\n%s\n\n(Enter to confirm, Esc to cancel)",
			m.Monitor,
			m.CustomInput.View(),
		)
	}

	s := fmt.Sprintf("Select Scale for %s\n\n", m.Monitor)

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true) // Pink/Orange
	normalStyle := lipgloss.NewStyle().PaddingLeft(2)
	currentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")) // Green

	for i, scale := range m.Scales {
		cursor := "  "
		if i == m.Selected {
			cursor = "▶ "
		}

		scaleStr := fmt.Sprintf("%.2fx", scale)
		indicator := ""
		if scale == 1.0 {
			indicator = " (native)"
		}
		if scale == m.Current {
			indicator = currentStyle.Render(" (current)")
		}

		line := fmt.Sprintf("%s%s", scaleStr, indicator)

		if i == m.Selected {
			s += selectedStyle.Render(cursor+line) + "\n"
		} else {
			s += normalStyle.Render(line) + "\n"
		}
	}

	sel := m.Scales[m.Selected]
	effW := int(float64(m.Width) / sel)
	effH := int(float64(m.Height) / sel)

	s += fmt.Sprintf("\nPhysical: %dx%d -> Effective: %dx%d\n", m.Width, m.Height, effW, effH)
	s += "\n[c] Custom  [1/2] Quick Select  [Esc] Cancel"

	return s
}
