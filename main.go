package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k0ryou/markdown-parser-go/generator"
	"github.com/k0ryou/markdown-parser-go/lexer"
	"github.com/k0ryou/markdown-parser-go/parser"
	"github.com/k0ryou/markdown-parser-go/token"
)

type Markdown struct {
	Content string
}

func main() {
	fmt.Println("start server")
	// POST http://localhost:8081/convertmd
	// curl -X POST -H "Content-Type: application/json" -d '{"Content": "{markdown text}"}' http://localhost:8081/convertmd
	http.HandleFunc("/convertmd", postMarkdownStringHandler)
	http.ListenAndServe(":8081", nil)
}

func postMarkdownStringHandler(w http.ResponseWriter, r *http.Request) {
	// json形式のMarkdownテキストを受け取る
	var md Markdown
	json.NewDecoder(r.Body).Decode(&md)

	// HTML変換後のテキストを出力
	fmt.Fprintln(w, convertToHTMLString(md.Content))
}

// MarkdownテキストをHTML形式に変換する
func convertToHTMLString(markdown string) string {
	// markdownを適切な区間ごとに分割する
	mdArray := lexer.Analize(markdown)
	var asts = [][]token.Token{}
	for _, md := range mdArray {
		// 一行ずつ抽象構文木に変換
		asts = append(asts, parser.Parse(md))
	}
	//抽象構文木からHTMLに変換
	htmlString := generator.Generate(asts)

	return htmlString
}
