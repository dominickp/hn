package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dominickp/go-hn-cli/client"
	"github.com/dominickp/go-hn-cli/util"
)

// TODO: implement viewport for scrolling https://github.com/charmbracelet/bubbletea/blob/master/examples/pager/main.go

const logfilePath = "logs/bubbletea.log"

type model struct {
	choices         []string // items on the to-do list
	topMenuResponse client.TopMenuResponse
	cursor          int // which to-do list item our cursor is pointing at
	// selected map[int]struct{} // which to-do items are selected
	err          error // an error to display, if any
	currentItem  int   // which item is currently selected
	currentTopic client.Item
}

func checkTopMenu() tea.Msg {
	topMenuResponse, err := client.GetTopMenuResponse(10)

	if err != nil {
		// There was an error making our request. Wrap the error we received
		// in a message and return it.
		return errMsg{err}
	}

	// We received a response from the server. Return the HTTP status code
	// as a message.
	return topMenuMsg(topMenuResponse)
}

func checkTopic(topicID int) tea.Msg {
	item, err := client.GetItemWithComments(topicID, 10)

	if err != nil {
		// There was an error making our request. Wrap the error we received
		// in a message and return it.
		return errMsg{err}
	}

	// We received a response from the server. Return the HTTP status code
	// as a message.
	return topicMsg(item)
}

type topMenuMsg client.TopMenuResponse
type topicMsg client.Item
type errMsg struct{ err error }

// Error implements error.
func (e errMsg) Error() string {
	panic("unimplemented")
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
	return checkTopMenu
}

func (m model) InitTopic() tea.Cmd {
	return func() tea.Msg {
		return checkTopic(m.currentItem)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case topMenuMsg:
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
		return m, tea.ClearScreen

	case topicMsg:
		// The server returned a topic response message. Save it to our model.
		m.currentTopic = client.Item(msg)
		fmt.Println("current topic: ", m.currentTopic.Title)
		// Set comments as choices
		choices := make([]string, len(m.currentTopic.Comments))
		for i, comment := range m.currentTopic.Comments {
			commentText := html.UnescapeString(comment.Text)
			commentText = strings.ReplaceAll(commentText, "<p>", "\n")
			commentText = strings.ReplaceAll(commentText, "</p>", "\n")
			choices[i] = html.UnescapeString(commentText)
		}
		m.choices = choices
		return m, tea.ClearScreen
	case errMsg:
		// There was an error. Note it in the model. And tell the runtime
		// we're done and want to quit.
		m.err = msg
		return m, tea.Quit

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
			item := m.topMenuResponse.Items[m.cursor]
			m.currentItem = item.Id
			m.cursor = 0

			return m, tea.Cmd(m.InitTopic())

		case "backspace":
			m.currentItem = 0
			m.currentTopic = client.Item{}
			m.cursor = 0
			return m, tea.Cmd(m.Init())
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {

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

	// Send the UI for rendering

	return s
}

func main() {
	if logfilePath != "" {
		if _, err := tea.LogToFile(logfilePath, "simple"); err != nil {
			log.Fatal(err)
		}
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		log.Fatal(err)
		os.Exit(1)
	}
}
