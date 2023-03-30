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
			args: args{markdown: "**hoge huga**"},
			want: "<strong>hoge huga</strong>",
		},
		{
			name: "strong_bad1",
			args: args{markdown: "** hoge huga**"},
			want: "** hoge huga**",
		},
		{
			name: "strong_bad2",
			args: args{markdown: `a**"hoge"**`},
			want: `a**"hoge"**`,
		},
		{
			name: "strong_ok2",
			args: args{markdown: "hoge**huga**"},
			want: "hoge<strong>huga</strong>",
		},
		{
			name: "strong_under_ok1",
			args: args{markdown: "__hoge huga__"},
			want: "<strong>hoge huga</strong>",
		},
		{
			name: "strong_under_bad1",
			args: args{markdown: "__ hoge huga__"},
			want: "__ hoge huga__",
		},
		{
			name: "strong_under_bad2",
			args: args{markdown: "__\nhoge huga__"},
			want: "__hoge huga__",
		},
		{
			name: "strong_under_bad3",
			args: args{markdown: `a__"hoge"__`},
			want: `a__"hoge"__`,
		},
		{
			name: "strong_under_bad4",
			args: args{markdown: "hoge__huga__"},
			want: "hoge__huga__",
		},
		{
			name: "strong_under_bad5",
			args: args{markdown: "5__6__78"},
			want: "5__6__78",
		},
		{
			name: "strong_under_bad6",
			args: args{markdown: "пристаням__стремятся__"},
			want: "пристаням__стремятся__",
		},
		{
			name: "strong_under_ok2",
			args: args{markdown: "__hoge, __huga__, baz__"},
			want: "<strong>hoge, <strong>huga</strong>, baz</strong>",
		},
		{
			name: "strong_under_ok3",
			args: args{markdown: "hoge-__(huga)__"},
			want: "hoge-<strong>(huga)</strong>",
		},
		{
			name: "list_1",
			args: args{markdown: "- hoge\n-hoge\n1. hoge\n1.2."},
			want: "<ul><li>hoge</li></ul>-hoge<ol><li>hoge</li></ol>1.2.",
		},
		{
			name: "list_2",
			args: args{markdown: "1. hogege- hoge\n- hoge\n- hoge\n - hoge\n1. hoge"},
			want: "<ol><li>hogege- hoge</li></ol><ul><li>hoge</li><li>hoge</li><li>hoge</li></ul><ol><li>hoge</li></ol>",
		},
		{
			name: "list_3",
			args: args{markdown: "1. hogege- hoge\n- hoge\n\n\n\n\n- hoge\n       - hoge\n1. hoge"},
			want: "<ol><li>hogege- hoge</li></ol><ul><li>hoge</li><li>hoge</li><li>hoge</li></ul><ol><li>hoge</li></ol>",
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
