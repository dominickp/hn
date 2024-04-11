package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dominickp/hn/client"
	"github.com/dominickp/hn/util"
)

type model struct {
	choices         []string // items on the to-do list
	topMenuResponse client.TopMenuResponse
	cursor          int   // which to-do list item our cursor is pointing at
	err             error // an error to display, if any
	currentItem     int   // which item is currently selected
	currentTopic    client.Item
	ready           bool
	viewport        viewport.Model
	pageSize        int
	currentPage     int
}

func initialModel() model {
	return model{
		choices:     []string{},
		pageSize:    15,
		currentPage: 1,
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		// Don't draw the top menu until we have the viewport size ready
		if m.ready {
			return checkTopMenu(m.pageSize, m.currentPage) // Get the top 500 stories and save to our cache
		}
		return checkNothing()
	}
}

func (m model) RedrawPage() tea.Cmd {
	return func() tea.Msg {
		if len(m.topMenuResponse.Items) > 0 {
			return checkTopMenuPage(m.topMenuResponse, m.pageSize, m.currentPage) // Get the top 500 stories and save to our cache
		}
		return checkNothing()
	}
}

func (m model) InitTopic() tea.Cmd {
	return func() tea.Msg {
		return checkTopic(m.currentItem)
	}
}

func getTopMenuCurrentPageChoices(m model) []string {
	choices := make([]string, m.pageSize)
	style := lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("8"))
	start := (m.currentPage - 1) * m.pageSize
	end := min(start+m.pageSize, len(m.topMenuResponse.Items))
	for i, item := range m.topMenuResponse.Items[start:end] {
		choices[i] = fmt.Sprintf("%s %s", style.Render(util.PadRight(strconv.Itoa(item.Score), 4)), item.Title)
	}
	return choices

}

// Update is called when "things happen." Its job is to look at what has happened and return an updated model in
// response. It can also return a Cmd to make more things happen.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {

	case topMenuMsg:
		// The server returned a top menu response message. Save it to our model.
		m.topMenuResponse = client.TopMenuResponse(msg)
		m.choices = getTopMenuCurrentPageChoices(m)
		m.viewport.SetContent(getContent(m))
		return m, tea.ClearScreen

	case checkTopMenuPageMsg:
		m.choices = getTopMenuCurrentPageChoices(m)
		m.viewport.SetContent(getContent(m))
		return m, tea.ClearScreen
	case topicMsg:
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
			commentText := util.HtmlToText(comment.Text)
			replies := ""
			if len(comment.Kids) > 0 {
				replies = fmt.Sprintf(" (%d replies)", len(comment.Kids))
			}
			choices[i] = fmt.Sprintf("%s\n%s", authorStyle.Render(comment.By+replies), textStyle.Render(commentText))
		}
		m.choices = choices
		m.viewport.SetContent(getContent(m))
		return m, tea.ClearScreen

	case errMsg:
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
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(getContent(m))
			m.ready = true

			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
			// Updatee the page size to call for more items per page if we can fit them
			m.pageSize = m.viewport.Height - 1
			return m, tea.Cmd(m.Init())
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

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
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "right":
			if m.currentTopic.Id == 0 {
				m.currentPage++
				return m, tea.Cmd(m.RedrawPage())
			}
		case "left":
			if m.currentTopic.Id == 0 && m.currentPage > 1 {
				m.currentPage--
				return m, tea.Cmd(m.RedrawPage())
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":

			// Find the item that the cursor is pointing at
			var item client.Item
			if m.currentTopic.Id != 0 {
				item = m.currentTopic.Comments[m.cursor]
			} else {
				itemIndex := (m.currentPage-1)*m.pageSize + m.cursor
				item = m.topMenuResponse.Items[itemIndex]
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
			return m, tea.Cmd(m.RedrawPage())

		case "f5":
			// Clear cache of top menu items
			m.topMenuResponse = client.TopMenuResponse{}
			// Go back to page 1
			m.currentPage = 1
			return m, tea.Cmd(m.Init())

		}

	}

	m.viewport.SetContent(getContent(m))

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// getContent returns the content to be displayed in the viewport.
func getContent(m model) string {

	var s string = ""

	if m.currentTopic.Id != 0 {
		// Render topic view
		if m.currentTopic.Title != "" {
			s += fmt.Sprintf("%s\n", util.TitleStyle.Render(m.currentTopic.Title))
		}

		if m.currentTopic.Text != "" {
			s += fmt.Sprintf("%s\n", util.TextStyle.Render(util.HtmlToText(m.currentTopic.Text)))
		}
		if m.currentTopic.Url != "" {
			s += fmt.Sprintf("→ %s\n", util.LinkStyle.Render(m.currentTopic.Url))
		}
		s += "\n"
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
	return s
}

// View is called when the program wants to render the UI. It returns a string.
func (m model) View() string {
	// Send the UI for rendering
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

// headerView returns the header view for the paginated viewport.
func (m model) headerView() string {
	title := util.TitleBoxStyle.Render("Hacker News")
	line := strings.Repeat("─", util.Max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

// footerView returns the footer view for the paginated viewport.
func (m model) footerView() string {
	navMessage := fmt.Sprintf("Page %d. Press q to quit, ←/→ to paginate, F5 to refresh.", m.currentPage)
	if m.currentTopic.Id != 0 {
		// Topic view
		navMessage = "Press q to quit, backspace to go back."
	}
	navHelpLine := fmt.Sprintf("─── %s ", navMessage)

	info := util.InfoBoxStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := navHelpLine + strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info+navHelpLine)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
