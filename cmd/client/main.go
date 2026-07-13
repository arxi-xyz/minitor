package main

import (
	"flag"
	"fmt"
	"os"

	"minitor/config"
	"minitor/view/terminal"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	configPath := flag.String("config", "", "path to JSON config file")
	url := flag.String("url", "", "websocket server url")
	flag.Parse()

	cfg, err := config.Load(config.LoadOptions{
		Path:  *configPath,
		WSURL: *url,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(terminal.NewModel(cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
