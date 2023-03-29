package main

import (
	"fmt"
	"markdown-parser-go/generator"
	"markdown-parser-go/lexer"
	"markdown-parser-go/parser"
	"markdown-parser-go/token"
)

func main() {
	inputStr := "text text\n- u -u 1. 2222. 2\n- 1. **strong1.**"
	fmt.Println(convertToHTMLString(inputStr))
}

func convertToHTMLString(markdown string) string {
	mdArray := lexer.Analize(markdown)
	var asts = [][]token.Token{}
	for _, md := range mdArray {
		asts = append(asts, parser.Parse(md))
	}
	htmlString := generator.Generate(asts)

	return htmlString
}
