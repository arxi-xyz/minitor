package main

import "minitor/transport"

func main() {
	// p := tea.NewProgram(view.InitialModel())
	// if _, err := p.Run(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	// }

	transport.NewServer().Run()
}
