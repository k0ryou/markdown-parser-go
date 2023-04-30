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
		want string
	}{
		{
			name: "text",
			args: args{markdown: "text**strong**\n- hoge\n- li\n1. ge"},
			want: "text<strong>strong</strong>\n<ul><li>hoge</li><li>li</li></ul><ol><li>ge</li></ol>",
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
			want: `a**&#34;hoge&#34;**`,
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
			want: "__\nhoge huga__",
		},
		{
			name: "strong_under_bad3",
			args: args{markdown: `a__"hoge"__`},
			want: `a__&#34;hoge&#34;__`,
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
			want: "<ul><li>hoge</li></ul>-hoge\n<ol><li>hoge</li></ol>1.2.",
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
		{
			name: "header_1",
			args: args{markdown: "# hoge - list"},
			want: "<h1>hoge - list</h1>",
		},
		{
			name: "header_2",
			args: args{markdown: "## hoge ## hogehoge"},
			want: "<h2>hoge ## hogehoge</h2>",
		},
		{
			name: "header_3",
			args: args{markdown: "- hoge # hoge"},
			want: "<ul><li>hoge # hoge</li></ul>",
		},
		{
			name: "header_4",
			args: args{markdown: "- ### hoge"},
			want: "<ul><li><h3>hoge</h3></li></ul>",
		},
		{
			name: "header_5",
			args: args{markdown: "**# hoge**"},
			want: "<strong># hoge</strong>",
		},
		{
			name: "header_6",
			args: args{markdown: "# **hoge**"},
			want: "<h1><strong>hoge</strong></h1>",
		},
		{
			name: "header_7",
			args: args{markdown: "####### hoge"},
			want: "####### hoge",
		},
		{
			name: "header_8",
			args: args{markdown: "1. # hoge"},
			want: "<ol><li><h1>hoge</h1></li></ol>",
		},
		{
			name: "a_1",
			args: args{markdown: "[hoge](hoge.com)"},
			want: "<a href='hoge.com'>hoge</a>",
		},
		{
			name: "a_2",
			args: args{markdown: "- [hoge](hoge.com)"},
			want: "<ul><li><a href='hoge.com'>hoge</a></li></ul>",
		},
		{
			name: "a_3",
			args: args{markdown: "**[hoge](hoge.com)**"},
			want: "<strong><a href='hoge.com'>hoge</a></strong>",
		},
		{
			name: "a_4",
			args: args{markdown: "1. **[hoge](hoge.com)**"},
			want: "<ol><li><strong><a href='hoge.com'>hoge</a></strong></li></ol>",
		},
		{
			name: "a_5",
			args: args{markdown: "[#h1](hoge.com)"},
			want: "<a href='hoge.com'>#h1</a>",
		},
		{
			name: "a_6",
			args: args{markdown: "[hoge](#hoge)"},
			want: "<a href='#hoge'>hoge</a>",
		},
		{
			name: "a_7",
			args: args{markdown: "hoge\n[hoge](hoge.com)text\n**strong**"},
			want: "hoge\n<a href='hoge.com'>hoge</a>text\n<strong>strong</strong>",
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
