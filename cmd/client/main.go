package main

import (
	"flag"
	"fmt"
	"os"

	"minitor/view/terminal"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	url := flag.String("url", "ws://127.0.0.1:8080/ws", "websocket server url")
	flag.Parse()

	p := tea.NewProgram(terminal.InitialModelWithURL(*url), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
