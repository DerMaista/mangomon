package tui

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"mangomon/config"

	"github.com/charmbracelet/lipgloss"
)

var (
	monitorBoxActive = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("42")).
				Foreground(lipgloss.Color("42"))

	monitorBoxInactive = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("244")).
				Foreground(lipgloss.Color("244"))

	monitorBoxSelected = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("214")).
				Foreground(lipgloss.Color("214"))

	monitorBoxMirror = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("205")).
				Foreground(lipgloss.Color("205"))
)

type GridModel struct {
	Rules         *map[string]config.MonitorRule
	SelectedID    string
	GridSize      int
	Width, Height int
}

func NewGridModel(rules *map[string]config.MonitorRule) GridModel {
	return GridModel{
		Rules:    rules,
		GridSize: 1,
	}
}

func (g GridModel) Bounds() (minX, minY, maxX, maxY int) {
	minX, minY = math.MaxInt, math.MaxInt
	maxX, maxY = math.MinInt, math.MinInt

	if len(*g.Rules) == 0 {
		return 0, 0, 1920, 1080
	}

	for _, r := range *g.Rules {
		if r.X < minX {
			minX = r.X
		}
		if r.Y < minY {
			minY = r.Y
		}
		if r.X+r.Width > maxX {
			maxX = r.X + r.Width
		}
		if r.Y+r.Height > maxY {
			maxY = r.Y + r.Height
		}
	}
	return
}

func (g GridModel) Render(termWidth, termHeight int) string {
	minX, minY, maxX, maxY := g.Bounds()
	totalW := maxX - minX
	totalH := maxY - minY

	if totalW <= 0 {
		totalW = 1920
	}
	if totalH <= 0 {
		totalH = 1080
	}

	paddingX := 3000
	paddingY := 3000

	viewMinX := minX - paddingX/2
	viewMinY := minY - paddingY/2
	viewW := totalW + paddingX
	viewH := totalH + paddingY

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(0, 1).
		Render(fmt.Sprintf("Grid: %d px", g.GridSize))

	renderHeight := termHeight - 2
	if renderHeight < 10 {
		renderHeight = 10
	}

	desktop := make([][]rune, renderHeight)
	for i := range desktop {
		desktop[i] = make([]rune, termWidth)
		for j := range desktop[i] {
			desktop[i][j] = ' '
		}
	}

	scaleX := float64(termWidth) / float64(viewW)
	scaleY := float64(renderHeight) / float64(viewH)

	termAspect := 2.2
	if scaleX > scaleY*termAspect {
		scaleX = scaleY * termAspect
	} else {
		scaleY = scaleX / termAspect
	}

	offsetX := (termWidth - int(float64(viewW)*scaleX)) / 2
	offsetY := (renderHeight - int(float64(viewH)*scaleY)) / 2

	worldToTerm := func(wx, wy int) (int, int) {
		tx := offsetX + int(float64(wx-viewMinX)*scaleX)
		ty := offsetY + int(float64(wy-viewMinY)*scaleY)
		return tx, ty
	}

	var ids []string
	for id := range *g.Rules {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		if ids[i] == g.SelectedID {
			return false
		}
		if ids[j] == g.SelectedID {
			return true
		}
		return ids[i] < ids[j]
	})

	for _, id := range ids {
		r := (*g.Rules)[id]

		x1, y1 := worldToTerm(r.X, r.Y)
		x2, y2 := worldToTerm(r.X+r.Width, r.Y+r.Height)

		if x2-x1 < 6 {
			x2 = x1 + 6
		}
		if y2-y1 < 4 {
			y2 = y1 + 4
		}

		if x1 < 0 {
			x1 = 0
		}
		if y1 < 0 {
			y1 = 0
		}
		if x2 >= termWidth {
			x2 = termWidth - 1
		}
		if y2 >= renderHeight {
			y2 = renderHeight - 1
		}

		style := monitorBoxInactive
		isActive := true // TODO: check if enabled
		if id == g.SelectedID {
			style = monitorBoxSelected
		} else if isActive {
			style = monitorBoxActive
		}

		box := getBoxRunes(style)
		drawBox(desktop, x1, y1, x2, y2, box)

		status := "[ON]"
		if !isActive {
			status = "[OFF]"
		}
		nameLabel := fmt.Sprintf("%s %s", id, status)
		drawText(desktop, x1+1, y1+1, x2-1, nameLabel)

		// 2. Resolution
		resLabel := fmt.Sprintf("%dx%d@%.0fHz", r.Width, r.Height, r.RefreshRate)
		drawText(desktop, x1+1, y1+2, x2-1, resLabel)

		// 3. Scale
		if r.Scale < 0.99 || r.Scale > 1.01 {
			scaleLabel := fmt.Sprintf("x%.2f", r.Scale)
			drawText(desktop, x1+1, y1+3, x2-1, scaleLabel)
		}
	}

	// Convert to string
	var access strings.Builder
	for _, row := range desktop {
		access.WriteString(string(row) + "\n")
	}

	return header + "\n" + access.String()
}

// Drawing Helpers

type boxRunes struct {
	topLeft     rune
	topRight    rune
	bottomLeft  rune
	bottomRight rune
	horizontal  rune
	vertical    rune
}

func getBoxRunes(style lipgloss.Style) boxRunes {
	border := style.GetBorderStyle()
	// Fallback if no border style
	if border.TopLeft == "" {
		return boxRunes{'┌', '┐', '└', '┘', '─', '│'}
	}
	return boxRunes{
		topLeft:     []rune(border.TopLeft)[0],
		topRight:    []rune(border.TopRight)[0],
		bottomLeft:  []rune(border.BottomLeft)[0],
		bottomRight: []rune(border.BottomRight)[0],
		horizontal:  []rune(border.Top)[0],
		vertical:    []rune(border.Left)[0],
	}
}

func drawBox(desktop [][]rune, x1, y1, x2, y2 int, box boxRunes) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			var ch rune
			if y == y1 && x == x1 {
				ch = box.topLeft
			} else if y == y1 && x == x2 {
				ch = box.topRight
			} else if y == y2 && x == x1 {
				ch = box.bottomLeft
			} else if y == y2 && x == x2 {
				ch = box.bottomRight
			} else if y == y1 || y == y2 {
				ch = box.horizontal
			} else if x == x1 || x == x2 {
				ch = box.vertical
			} else {
				continue // Don't fill center yet
			}
			desktop[y][x] = ch
		}
	}
}

func drawText(desktop [][]rune, x, y, maxX int, text string) {
	if y >= len(desktop) || y < 0 {
		return
	}
	runes := []rune(text)
	for i, r := range runes {
		if x+i >= len(desktop[0]) || x+i > maxX {
			break
		}
		desktop[y][x+i] = r
	}
}

func splitLines(s string) []string {
	var lines []string
	cur := ""
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, cur)
			cur = ""
		} else {
			cur += string(r)
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

func (g *GridModel) MoveSelected(dx, dy int) {
	if rule, ok := (*g.Rules)[g.SelectedID]; ok {
		stepX := dx * g.GridSize
		stepY := dy * g.GridSize

		rule.X += stepX
		rule.Y += stepY
		(*g.Rules)[g.SelectedID] = rule
	}
}

func (g *GridModel) CycleGrid() {
	switch g.GridSize {
	case 1:
		g.GridSize = 8
	case 8:
		g.GridSize = 16
	case 16:
		g.GridSize = 32
	case 32:
		g.GridSize = 64
	default:
		g.GridSize = 1
	}
}
