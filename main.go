package main

import (
	"encoding/json"
	"fmt"
	"markdown-parser-go/generator"
	"markdown-parser-go/lexer"
	"markdown-parser-go/parser"
	"markdown-parser-go/token"
	"net/http"
)

type Markdown struct {
	Content string
}

func main() {
	http.HandleFunc("/convertmd", postMarkdownStringHandler)
	http.ListenAndServe(":8081", nil)
}

func postMarkdownStringHandler(w http.ResponseWriter, r *http.Request) {
	var md Markdown
	json.NewDecoder(r.Body).Decode(&md)

	fmt.Fprintln(w, convertToHTMLString(md.Content))
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
