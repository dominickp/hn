package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dominickp/hn/client"
)

func checkTopMenu(pageSize, page int) tea.Msg {
	topMenuResponse, err := client.GetTopMenuResponse(pageSize, page)

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

func checkNothing() tea.Msg {
	return nil
}

type topMenuMsg client.TopMenuResponse
type topicMsg client.Item
type errMsg struct{ err error }

// Error implements error.
func (e errMsg) Error() string {
	panic("unimplemented")
}
