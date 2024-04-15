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
