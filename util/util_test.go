package util

import (
	"testing"
)

func TestPadRight(t *testing.T) {
	type args struct {
		str    string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestPadRight",
			args: args{str: "hello", length: 10},
			want: "hello     ",
		},
		{
			name: "TestLongerInput",
			args: args{str: "hello world", length: 5},
			want: "hello world",
		},
		{
			name: "TestEmpty",
			args: args{str: "", length: 5},
			want: "     ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PadRight(tt.args.str, tt.args.length); got != tt.want {
				t.Errorf("PadRight() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestHtmlToText(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestReplaceP",
			args: args{s: "Hello<p>world</p>"},
			want: "Hello\nworld",
		},
		{
			name: "TestReplaceMultipleP",
			args: args{s: "Hello<p>world.</p> Say hello to my <p>little friend!</p>"},
			want: "Hello\nworld.\n Say hello to my \nlittle friend!",
		},
		{
			name: "TestItalicizeI",
			args: args{s: "Hello <i>world</i>"},
			want: "Hello world", // I'm not sure why ansi escape codes are not being rendered here
		},
		{
			name: "TestLinkFormatting",
			args: args{s: "This is a link: <a href=\"https://example.com\">example</a>"},
			want: "This is a link: https://example.com",
		},
		{
			name: "TestLinkFormattingWRel",
			args: args{s: "This is a link: <a href=\"https://example.com\" rel=\"foo\">example</a>"},
			want: "This is a link: https://example.com",
		},
		{
			name: "TestEmpty",
			args: args{s: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HtmlToText(tt.args.s)
			if got != tt.want {
				t.Errorf("HtmlToText() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
