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
			args: args{markdown: `a**"foo"**`},
			want: `a**"foo"**`,
		},
		{
			name: "strong_ok2",
			args: args{markdown: "foo**bar**"},
			want: "foo<strong>bar</strong>",
		},
		{
			name: "strong_under_ok1",
			args: args{markdown: "__foo bar__"},
			want: "<strong>foo bar</strong>",
		},
		{
			name: "strong_under_bad1",
			args: args{markdown: "__ foo bar__"},
			want: "__ foo bar__",
		},
		{
			name: "strong_under_bad2",
			args: args{markdown: "__\nfoo bar__"},
			want: "__foo bar__",
		},
		{
			name: "strong_under_bad3",
			args: args{markdown: `a__"foo"__`},
			want: `a__"foo"__`,
		},
		{
			name: "strong_under_bad4",
			args: args{markdown: "foo__bar__"},
			want: "foo__bar__",
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
			args: args{markdown: "__foo, __bar__, baz__"},
			want: "<strong>foo, <strong>bar</strong>, baz</strong>",
		},
		{
			name: "strong_under_ok3",
			args: args{markdown: "foo-__(bar)__"},
			want: "foo-<strong>(bar)</strong>",
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
