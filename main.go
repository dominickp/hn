package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const logfilePath = "logs/bubbletea.log"
const useHighPerformanceRenderer = false

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	if logfilePath != "" {
		if _, err := tea.LogToFile(logfilePath, "simple"); err != nil {
			log.Fatal(err)
		}
	}

	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		log.Fatal(err)
		os.Exit(1)
	}
}
