package lexer

import (
	"reflect"
	"testing"
)

func TestAnalize(t *testing.T) {
	type args struct {
		markdown string
	}
	tests := []struct {
		name string
		args args
		// listの要素はそれぞれに改行文字を付与する
		want []string
	}{
		{
			name: "normal",
			args: args{markdown: "text text\n- ul\n1. ol\n**strong**"},
			want: []string{"text text", "- ul\n", "1. ol\n", "**strong**"},
		},
		{
			name: "list",
			args: args{markdown: "text text\n- ul\n- ul2\n1. **strong_ol**\n2. ol2\n- ul3"},
			want: []string{"text text", "- ul\n- ul2\n", "1. **strong_ol**\n2. ol2\n", "- ul3\n"},
		},
		{
			name: "hack",
			args: args{markdown: "text text\n- u -u 1. 2222. 2\n- 1. **strong1.**"},
			want: []string{"text text", "- u -u 1. 2222. 2\n- 1. **strong1.**\n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Analize(tt.args.markdown)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Analize() = %v, want %v", got, tt.want)
			}
		})
	}
}
