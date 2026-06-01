package main

import (
	"reflect"
	"testing"
)

func TestSplitEnvList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{
			name: "空文字は nil",
			raw:  "",
			want: nil,
		},
		{
			name: "trim と重複除去を行う",
			raw:  " https://www.example.com ,https://api.example.com,https://www.example.com, ",
			want: []string{
				"https://www.example.com",
				"https://api.example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := splitEnvList(tt.raw); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("splitEnvList(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}
