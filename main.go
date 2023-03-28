package main

import (
	"fmt"
	"markdown-parser-go/generator"
	"markdown-parser-go/lexer"
	"markdown-parser-go/parser"
	"markdown-parser-go/token"
)

func main() {
	inputStr := "normal text\n\n- __boldlist1__\n- list2\n"
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
