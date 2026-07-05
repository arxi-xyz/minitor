package main

import (
	"fmt"
	"os"

	view "minitor/view/terminal"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(view.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}
