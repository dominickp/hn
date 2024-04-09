package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dominickp/go-hn-cli/client"
	"github.com/dominickp/go-hn-cli/messages"
	"github.com/dominickp/go-hn-cli/util"
)

// TODO: implement viewport for scrolling https://github.com/charmbracelet/bubbletea/blob/master/examples/pager/main.go
var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

const logfilePath = "logs/bubbletea.log"
const useHighPerformanceRenderer = false

type model struct {
	choices         []string // items on the to-do list
	topMenuResponse client.TopMenuResponse
	cursor          int // which to-do list item our cursor is pointing at
	// selected map[int]struct{} // which to-do items are selected
	err          error // an error to display, if any
	currentItem  int   // which item is currently selected
	currentTopic client.Item
	ready        bool
	viewport     viewport.Model
	content      string
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []string{},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		// selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return messages.CheckTopMenu
}

func (m model) InitTopic() tea.Cmd {
	return func() tea.Msg {
		return messages.CheckTopic(m.currentItem)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {

	case messages.TopMenuMsg:
		// The server returned a top menu response message. Save it to our model.
		m.topMenuResponse = client.TopMenuResponse(msg)
		choices := make([]string, len(msg.Items))
		style := lipgloss.NewStyle().
			Bold(false).
			Foreground(lipgloss.Color("8"))
		for i, item := range msg.Items {
			choices[i] = fmt.Sprintf("%s %s", style.Render(util.PadRight(strconv.Itoa(item.Score), 4)), item.Title)
		}
		m.choices = choices
		m.viewport.SetContent(getContent(m))
		return m, tea.ClearScreen

	case messages.TopicMsg:
		// The server returned a topic response message. Save it to our model.
		m.currentTopic = client.Item(msg)
		fmt.Println("current topic: ", m.currentTopic.Title)
		// Styles
		authorStyle := lipgloss.NewStyle().
			Bold(false).
			Foreground(lipgloss.Color("8"))
		textStyle := lipgloss.NewStyle().
			PaddingLeft(4).
			PaddingBottom(2)
		// Set comments as choices
		choices := make([]string, len(m.currentTopic.Comments))
		for i, comment := range m.currentTopic.Comments {
			// TODO
			// Handle italics
			// Handle quotes (line starts with >), color green
			// Handle links
			commentText := html.UnescapeString(comment.Text)
			commentText = strings.ReplaceAll(commentText, "<p>", "\n")
			commentText = strings.ReplaceAll(commentText, "</p>", "\n")
			replies := ""
			if len(comment.Kids) > 0 {
				replies = fmt.Sprintf(" (%d replies)", len(comment.Kids))
			}
			choices[i] = fmt.Sprintf("%s\n%s", authorStyle.Render(comment.By+replies), textStyle.Render(html.UnescapeString(commentText)))
		}
		m.choices = choices
		m.viewport.SetContent(getContent(m))
		return m, tea.ClearScreen
	case messages.ErrMsg:
		// There was an error. Note it in the model. And tell the runtime
		// we're done and want to quit.
		m.err = msg
		return m, tea.Quit

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(getContent(m))
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		// if useHighPerformanceRenderer {
		// 	// Render (or re-render) the whole viewport. Necessary both to
		// 	// initialize the viewport and when the window is resized.
		// 	//
		// 	// This is needed for high-performance rendering only.
		// 	cmds = append(cmds, viewport.Sync(m.viewport))
		// }

	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			// FIXME: only allow scrolling on top menu
			//			Until I figure out what I'm doing for the topic view
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":

			// Find the item that the cursor is pointing at
			var item client.Item
			if m.currentTopic.Id != 0 {
				item = m.currentTopic.Comments[m.cursor]
			} else {
				item = m.topMenuResponse.Items[m.cursor]
			}

			m.currentItem = item.Id
			m.cursor = 0
			m.viewport.GotoTop()
			return m, tea.Cmd(m.InitTopic())

		case "backspace":
			m.currentItem = 0
			m.currentTopic = client.Item{}
			m.cursor = 0
			m.viewport.GotoTop()
			return m, tea.Cmd(m.Init())
		}
		m.viewport.SetContent(getContent(m))
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func getContent(m model) string {

	var s string = ""

	if m.currentTopic.Id != 0 {
		// Render topic view
		s += fmt.Sprintf("Title: %s\n", m.currentTopic.Title)
		if m.currentTopic.Text != "" {
			s += fmt.Sprintf("Text: %s\n", m.currentTopic.Text)
		}
		if m.currentTopic.Url != "" {
			s += fmt.Sprintf("URL: %s\n", m.currentTopic.Url)
		}
		s += "\n"
	} else {
		// Render top menu view
		s += "HackerNews Top Topics:\n\n"
	}

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		style := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("5"))
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", style.Render(cursor), choice)
	}

	// The footer
	s += "\nPress q to quit, backspace to go back.\n"
	return s
}

func (m model) View() string {

	// var s string = ""

	// if m.currentTopic.Id != 0 {
	// 	// Render topic view
	// 	s += fmt.Sprintf("Title: %s\n", m.currentTopic.Title)
	// 	if m.currentTopic.Text != "" {
	// 		s += fmt.Sprintf("Text: %s\n", m.currentTopic.Text)
	// 	}
	// 	if m.currentTopic.Url != "" {
	// 		s += fmt.Sprintf("URL: %s\n", m.currentTopic.Url)
	// 	}
	// 	s += "\n"
	// } else {
	// 	// Render top menu view
	// 	s += "HackerNews Top Topics:\n\n"
	// }

	// // Iterate over our choices
	// for i, choice := range m.choices {

	// 	// Is the cursor pointing at this choice?
	// 	style := lipgloss.NewStyle().
	// 		Bold(true).
	// 		Foreground(lipgloss.Color("5"))
	// 	cursor := " " // no cursor
	// 	if m.cursor == i {
	// 		cursor = ">" // cursor!
	// 	}

	// 	// Render the row
	// 	s += fmt.Sprintf("%s %s\n", style.Render(cursor), choice)
	// }

	// // The footer
	// s += "\nPress q to quit, backspace to go back.\n"

	// Send the UI for rendering
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

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
