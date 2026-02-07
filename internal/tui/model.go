package tui

import (
	"fmt"
	"mangomon/config"
	"mangomon/internal/state"
	"mangomon/internal/system"
	"mangomon/internal/tui/tools"

	tea "github.com/charmbracelet/bubbletea"
)

type modelState int

const (
	stateGrid modelState = iota
	stateScale
	stateMode
	stateMirror
	stateTransform
	stateVRR
)

type Model struct {
	outputs []system.Output
	rules   map[string]config.MonitorRule
	state   modelState
	parser  *config.ConfigParser
	err     error

	// Grid state
	grid GridModel

	// Tools
	scalePicker     tools.ScalePickerModel
	modePicker      tools.ModePickerModel
	mirrorPicker    tools.MirrorPickerModel
	transformPicker tools.TransformPickerModel
	vrrPicker       tools.VRRPickerModel

	width, height int
}

func InitialModel(parser *config.ConfigParser) Model {
	outputs, err := system.GetOutputs()
	// Mock outputs if none found for testing
	if len(outputs) == 0 {
		// Fallback or just empty
	}

	rules, _ := parser.Parse()

	// Ensure every output has a rule
	for _, out := range outputs {
		if _, ok := rules[out.Name]; !ok {
			rules[out.Name] = config.MonitorRule{
				ID:                  out.Name,
				Scale:               1.0,
				Transform:           0,
				VariableRefreshRate: 0,
				X:                   0, Y: 0,
				Width: 1920, Height: 1080, RefreshRate: 60,
			}
		}
	}

	initialSelected := ""
	if len(outputs) > 0 {
		initialSelected = outputs[0].Name
	}

	grid := NewGridModel(&rules)
	grid.SelectedID = initialSelected

	// Load app state (GridSize only)
	if appState, err := state.Load(); err == nil {
		grid.GridSize = appState.GridSize
		if grid.GridSize == 0 {
			grid.GridSize = 1
		}
	}

	return Model{
		outputs: outputs,
		rules:   rules,
		parser:  parser,
		err:     err,
		grid:    grid,
		state:   stateGrid,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	// Sub-model returns
	case tools.ScaleSelectedMsg:
		if rule, ok := m.rules[m.grid.SelectedID]; ok {
			rule.Scale = msg.Scale
			m.rules[m.grid.SelectedID] = rule
		}
		m.state = stateGrid
		m.grid.Rules = &m.rules // Refresh grid pointer just in case
		return m, nil

	case tools.ScaleCancelledMsg:
		m.state = stateGrid
		return m, nil

	case tools.ModeSelectedMsg:
		if rule, ok := m.rules[m.grid.SelectedID]; ok {
			rule.Width = msg.Mode.Width
			rule.Height = msg.Mode.Height
			rule.RefreshRate = msg.Mode.Rate
			m.rules[m.grid.SelectedID] = rule
		}
		m.state = stateGrid
		m.grid.Rules = &m.rules
		return m, nil

	case tools.ModeCancelledMsg:
		m.state = stateGrid
		return m, nil

	case tools.MirrorSelectedMsg:
		// Align position and resolution to target
		if rule, ok := m.rules[m.grid.SelectedID]; ok {
			if target, ok := m.rules[msg.TargetID]; ok {
				rule.X = target.X
				rule.Y = target.Y
				rule.Width = target.Width
				rule.Height = target.Height
			}
			m.rules[m.grid.SelectedID] = rule
		}
		m.state = stateGrid
		return m, nil

	case tools.MirrorCancelledMsg:
		m.state = stateGrid
		return m, nil

	case tools.TransformSelectedMsg:
		if rule, ok := m.rules[m.grid.SelectedID]; ok {
			rule.Transform = msg.Transform
			m.rules[m.grid.SelectedID] = rule
		}
		m.state = stateGrid
		m.grid.Rules = &m.rules
		return m, nil

	case tools.TransformCancelledMsg:
		m.state = stateGrid
		return m, nil

	case tools.VRRSelectedMsg:
		if rule, ok := m.rules[m.grid.SelectedID]; ok {
			rule.VariableRefreshRate = msg.VRR
			m.rules[m.grid.SelectedID] = rule
		}
		m.state = stateGrid
		m.grid.Rules = &m.rules
		return m, nil

	case tools.VRRCancelledMsg:
		m.state = stateGrid
		return m, nil

	}

	// Delegate based on state
	switch m.state {
	case stateGrid:
		return m.updateGrid(msg)
	case stateScale:
		newScale, cmd := m.scalePicker.Update(msg)
		m.scalePicker = newScale.(tools.ScalePickerModel)
		return m, cmd
	case stateMode:
		newModel, cmd := m.modePicker.Update(msg)
		m.modePicker = newModel.(tools.ModePickerModel)
		return m, cmd
	case stateMirror:
		newModel, cmd := m.mirrorPicker.Update(msg)
		m.mirrorPicker = newModel.(tools.MirrorPickerModel)
		return m, cmd
	case stateTransform:
		newModel, cmd := m.transformPicker.Update(msg)
		m.transformPicker = newModel.(tools.TransformPickerModel)
		return m, cmd
	case stateVRR:
		newModel, cmd := m.vrrPicker.Update(msg)
		m.vrrPicker = newModel.(tools.VRRPickerModel)
		return m, cmd
	}

	return m, nil
}

func (m Model) updateGrid(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			currIdx := -1
			for i, o := range m.outputs {
				if o.Name == m.grid.SelectedID {
					currIdx = i
					break
				}
			}
			if currIdx != -1 {
				nextIdx := (currIdx + 1) % len(m.outputs)
				m.grid.SelectedID = m.outputs[nextIdx].Name
			}
			if len(m.outputs) == 0 {
				// Handle no output case
			} else if currIdx == -1 {
				m.grid.SelectedID = m.outputs[0].Name
			}

		case "up", "k":
			m.grid.MoveSelected(0, -1)
		case "down", "j":
			m.grid.MoveSelected(0, 1)
		case "left", "h":
			m.grid.MoveSelected(-1, 0)
		case "right", "l":
			m.grid.MoveSelected(1, 0)

		case "shift+up", "K":
			m.grid.MoveSelected(0, -10)
		case "shift+down", "J":
			m.grid.MoveSelected(0, 10)
		case "shift+left", "H":
			m.grid.MoveSelected(-10, 0)
		case "shift+right", "L": // Initial confusion: L usually is right. But Shift+L is Shift+Right.
			// Wait, "L" is also Snap Toggle in my plan.
			// Hyprmon uses 'l' for right (vim key). 'L' (Shift+l) for fast right.
			// Does it use a different key for Snap?
			// "Cycle snap mode: s" (Wait, s is Save?)
			// Let's check Hyprmon keys.
			// Hyprmon: s=save? No, "enter"=save?
			// "s" -> "Toggle snap".
			// "Enter" -> "Apply/Save".
			// "r" -> "Scale". "m" -> "Mirror/Mode".
			// My implementation plan said:
			// S = Save, L = Snap.
			// Vim keys: h,j,k,l. Shift+L = Fast Right.
			// So "l" cannot be snap if "l" is right.
			// I will use 'z' for Snap for now to avoid conflict with 'l'.
			// Actually, I can check if key is 'l' (lower) vs 'L' (upper).
			// But Shift+l gives 'L'.
			// If I map 'L' to Fast Right, I can't use it for Snap.
			// I'll use 'p' for Snap? No, P is Profile.
			// Let's use 'n' (Re-snap/SsNap)?
			// I'll use 'G' for Grid Size (Shift+g). 'g' is home?
			// 'g' is home in vim.
			// I'll use 'G' for Grid Cycle.
			// I'll use 'b' for Snap (Border/Both)?
			// Let's use 'A' (Align/Snap).
			m.grid.MoveSelected(10, 0)

		case "G", "g":
			m.grid.CycleGrid()

		case "R", "r": // Open Scale Picker
			if rule, ok := m.rules[m.grid.SelectedID]; ok {
				m.state = stateScale
				m.scalePicker = tools.NewScalePicker(rule.ID, rule.Scale, rule.Width, rule.Height)
			}

		case "F", "f": // Open Mode Picker
			if rule, ok := m.rules[m.grid.SelectedID]; ok {
				sysModes, _ := system.GetModes(rule.ID) // Ignore error, returns default list if fail
				var toolModes []tools.Mode
				for _, sm := range sysModes {
					toolModes = append(toolModes, tools.Mode{Width: sm.Width, Height: sm.Height, Rate: sm.Rate})
				}

				m.state = stateMode
				m.modePicker = tools.NewModePicker(rule.ID, tools.Mode{Width: rule.Width, Height: rule.Height, Rate: rule.RefreshRate}, toolModes)
			}

		case "M", "m":
			var allNames []string
			for _, o := range m.outputs {
				allNames = append(allNames, o.Name)
			}
			m.state = stateMirror
			m.mirrorPicker = tools.NewMirrorPicker(m.grid.SelectedID, allNames)

		case "T", "t":
			if rule, ok := m.rules[m.grid.SelectedID]; ok {
				m.state = stateTransform
				m.transformPicker = tools.NewTransformPicker(rule.ID, rule.Transform)
			}

		case "V", "v": // Open VRR Picker
			if rule, ok := m.rules[m.grid.SelectedID]; ok {
				m.state = stateVRR
				m.vrrPicker = tools.NewVRRPicker(rule.ID, rule.VariableRefreshRate)
			}

		case "S", "s": // Save
			// Save app state (GridSize only)
			appState := state.AppState{
				GridSize: m.grid.GridSize,
			}
			if err := state.Save(appState); err != nil {
				m.err = err
			}

			var rulesToSave []config.MonitorRule
			for _, r := range m.rules {
				rulesToSave = append(rulesToSave, r)
			}
			err := m.parser.Save(rulesToSave)
			if err != nil {
				m.err = err
			}
			// Maybe show a "Saved!" message?
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateGrid:
		return m.viewGrid()
	case stateScale:
		return m.scalePicker.View()
	case stateMode:
		return m.modePicker.View()
	case stateMirror:
		return m.mirrorPicker.View()
	case stateTransform:
		return m.transformPicker.View()
	case stateVRR:
		return m.vrrPicker.View()
	}
	return ""
}

func (m Model) viewGrid() string {
	h := m.height - 4
	if h < 10 {
		h = 10
	}

	content := m.grid.Render(m.width, h)

	footer := "[Tab] Cycle  [Arrows] Move  [G] Grid  [R] Scale  [F] Mode  [T] Transform  [V] VRR  [M] Mirror  [S] Save  [Q] Quit"
	if m.err != nil {
		footer = fmt.Sprintf("Error: %v", m.err)
	}

	return fmt.Sprintf("MangoWC Spatial Config\n%s\n%s", content, footer)
}
