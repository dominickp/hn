package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dominickp/hn/client"
)

func Test_getTopMenuCurrentPageChoices(t *testing.T) {
	type args struct {
		m model
	}

	topMenuResponse := client.TopMenuResponse{}
	numItems := 25
	for i := 1; i < numItems; i++ {
		item := client.Item{Id: i, Title: fmt.Sprintf("item %d", i), Score: 33}
		topMenuResponse.Items = append(topMenuResponse.Items, item)
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "TestGetTopMenuCurrentPageChoices",
			args: args{m: model{currentPage: 1, pageSize: 5, topMenuResponse: topMenuResponse}},
			want: []string{"33   item 1", "33   item 2", "33   item 3", "33   item 4", "33   item 5"},
		},
		{
			name: "TestGetTopMenuCurrentPageChoicesPage2",
			args: args{m: model{currentPage: 2, pageSize: 5, topMenuResponse: topMenuResponse}},
			want: []string{"33   item 6", "33   item 7", "33   item 8", "33   item 9", "33   item 10"},
		},
		{
			name: "TestGetTopMenuCurrentPageLongerPage",
			args: args{m: model{currentPage: 1, pageSize: 10, topMenuResponse: topMenuResponse}},
			want: []string{
				"33   item 1", "33   item 2", "33   item 3", "33   item 4", "33   item 5",
				"33   item 6", "33   item 7", "33   item 8", "33   item 9", "33   item 10",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTopMenuCurrentPageChoices(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTopMenuCurrentPageChoices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_model_getCurrentTopic(t *testing.T) {
	tests := []struct {
		name string
		m    model
		want *client.Item
	}{
		{
			name: "TestGetCurrentTopic",
			m: model{topicHistoryStack: []client.Item{
				{Id: 1, Title: "item 1"},
				{Id: 1, Title: "item 2"},
				{Id: 1, Title: "item 3"},
			}},
			want: &client.Item{Id: 1, Title: "item 3"},
		},
		{
			name: "TestGetCurrentTopicEmptyStack",
			m:    model{topicHistoryStack: []client.Item{}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.getCurrentTopic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("model.getCurrentTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getContent(t *testing.T) {
	type args struct {
		m model
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetContentWithCurrentTopic",
			args: args{m: model{
				topicHistoryStack: []client.Item{{
					Id: 1, Title: "item 1", By: "Joe", Kids: []int{1, 2, 3}, Url: "http://example.com", Text: "foo...",
				}},
				choices: []string{"33   item 1", "33   item 2"},
			}},
			want: `item 1
By Joe (3 comments)
    foo...
→ http://example.com

> 33   item 1
  33   item 2
`,
		},
		{
			name: "TestGetContentTopMenu",
			args: args{m: model{
				topicHistoryStack: []client.Item{},
				choices:           []string{"33   item 1", "33   item 2", "33   item 3", "33   item 4", "33   item 5"},
			}},
			want: `> 33   item 1
  33   item 2
  33   item 3
  33   item 4
  33   item 5
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getContent(tt.args.m); got != tt.want {
				t.Errorf("getContent() = \n'%v', want \n'%v'", got, tt.want)
			}
		})
	}
}
