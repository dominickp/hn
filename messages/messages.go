package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dominickp/go-hn-cli/client"
)

func CheckTopMenu() tea.Msg {
	topMenuResponse, err := client.GetTopMenuResponse(15)

	if err != nil {
		// There was an error making our request. Wrap the error we received
		// in a message and return it.
		return ErrMsg{err}
	}

	// We received a response from the server. Return the HTTP status code
	// as a message.
	return TopMenuMsg(topMenuResponse)
}

func CheckTopic(topicID int) tea.Msg {
	item, err := client.GetItemWithComments(topicID, 10)

	if err != nil {
		// There was an error making our request. Wrap the error we received
		// in a message and return it.
		return ErrMsg{err}
	}

	// We received a response from the server. Return the HTTP status code
	// as a message.
	return TopicMsg(item)
}

type TopMenuMsg client.TopMenuResponse
type TopicMsg client.Item
type ErrMsg struct{ err error }

// Error implements error.
func (e ErrMsg) Error() string {
	panic("unimplemented")
}
