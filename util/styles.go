package util

import "github.com/charmbracelet/lipgloss"

var (
	TitleBoxStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	InfoBoxStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return TitleBoxStyle.Copy().BorderStyle(b)
	}()

	TextStyle   = lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("8"))
	LinkStyle   = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("1"))
	ItalicStyle = lipgloss.NewStyle().Italic(true)
	QuoteStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	TitleStyle  = lipgloss.NewStyle().Bold(true)
)
