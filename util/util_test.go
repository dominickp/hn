package util

import "testing"

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
			args: args{
				str:    "hello",
				length: 10,
			},
			want: "hello     ",
		},
		{
			name: "TestLongerInput",
			args: args{
				str:    "hello world",
				length: 5,
			},
			want: "hello world ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PadRight(tt.args.str, tt.args.length); got != tt.want {
				t.Errorf("PadRight() = %v, want %v", got, tt.want)
			}
		})
	}
}
