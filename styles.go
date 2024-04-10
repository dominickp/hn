package main

import "github.com/charmbracelet/lipgloss"

var (
	titleBoxStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoBoxStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleBoxStyle.Copy().BorderStyle(b)
	}()
	textStyle = lipgloss.NewStyle().
			Bold(false).
			Foreground(lipgloss.Color("8"))

	linkStyle   = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("1"))
	italicStyle = lipgloss.NewStyle().Italic(true)
	quoteStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	titleStyle  = lipgloss.NewStyle().Bold(true)
)
