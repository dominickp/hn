package main

// A simple program demonstrating the paginator component from the Bubbles
// component library.

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/dominickp/go-hn-cli/client"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	form         *huh.Form // huh.Form is just a tea.Model
	items        []client.Item
	selectedItem int
}

var selectedItem int

func NewModel() Model {
	//fmt.Println("msg: ", msg)
	topMenuResponse, _ := client.GetTopMenuResponse(10)
	items := topMenuResponse.Items

	huhOptions := make([]huh.Option[int], len(items))
	for i, item := range items {
		scoreKey := fmt.Sprintf("%d", item.Score)
		optionTitle := fmt.Sprintf("%s %s", scoreKey, item.Title) // Convert item.Score to string using %d.
		huhOptions[i] = huh.NewOption(optionTitle, item.Id)
	}

	return Model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[int]().
					Title("Top Topics").
					Options(huhOptions...).
					Key("selectedItem").Value(&selectedItem),
			),
		),
		items: items,
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// default:

	// }

	// ...

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	return m, cmd
}

func (m Model) View() string {
	if m.form.State == huh.StateCompleted {
		// selectedItem := m.form.GetString("selectedItem")
		return fmt.Sprintf("You selected: %d", selectedItem)
	}
	return m.form.View()
}

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
