package main

import (
	"testing"
)

func TestConvertToHTMLString(t *testing.T) {
	type args struct {
		markdown string
	}
	tests := []struct {
		name string
		args args
		// listの要素はそれぞれに改行文字を付与する
		want string
	}{
		{
			name: "text",
			args: args{markdown: "text**strong**\n- hoge\n- li\n1. ge"},
			want: "text<strong>strong</strong><ul><li>hoge</li><li>li</li></ul><ol><li>ge</li></ol>",
		},
		{
			name: "strong_ok1",
			args: args{markdown: "**foo bar**"},
			want: "<strong>foo bar</strong>",
		},
		{
			name: "strong_bad1",
			args: args{markdown: "** foo bar**"},
			want: "** foo bar**",
		},
		{
			name: "strong_bad2",
			args: args{markdown: 'a**"foo"**'},
			want: "a**&quot;foo&quot;**",
		},
		{
			name: "strong_ok2",
			args: args{markdown: "foo**bar**"},
			want: "foo<strong>bar</strong>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToHTMLString(tt.args.markdown)
			if got != tt.want {
				t.Errorf("converToHTMLString() = %v, want %v", got, tt.want)
			}
		})
	}
}
