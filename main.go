package main

import (
	"fmt"
	"os"

	"mangomon/config"
	"mangomon/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	parser, err := config.NewParser("")
	if err != nil {
		fmt.Printf("Error initializing parser: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(tui.InitialModel(parser))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
