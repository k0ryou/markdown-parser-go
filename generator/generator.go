package generator

import (
	"markdown-parser-go/token"
	"math"
	"regexp"
	"sort"
	"strings"
)

const (
	HTML_TAG_REGEXP = `<(.*?)>`
)

func Generate(asts [][]token.Token) string {
	htmlStrings := []string{}
	for _, ast := range asts {
		lineHtmlStrings := []string{}
		rearrangedAst := ast
		parIdxMap := map[int]int{}
		sort.Slice(rearrangedAst, func(i, j int) bool { return rearrangedAst[i].Id > rearrangedAst[j].Id })
		for idx, curToken := range rearrangedAst {
			parIdxMap[curToken.Id] = idx
		}
		for _, curToken := range rearrangedAst {
			if curToken.Parent.ElmType != token.ROOT {
				parIdx := parIdxMap[curToken.Parent.Id]
				parToken := rearrangedAst[parIdx]
				mergedToken := token.Token{
					Id:      parToken.Id,
					Parent:  parToken.Parent,
					ElmType: token.MERGED,
					Content: createMergedContent(curToken, parToken),
				}
				rearrangedAst[parIdx] = mergedToken
			} else {
				lineHtmlStrings = append(lineHtmlStrings, curToken.Content)
			}
		}

		loopLimit := int(math.Sqrt(float64(len(lineHtmlStrings))))
		for i := 0; i < loopLimit; i++ { // 逆順に並んでいるので反転する
			tmp := lineHtmlStrings[i]
			lineHtmlStrings[i] = lineHtmlStrings[len(lineHtmlStrings)-i-1]
			lineHtmlStrings[len(lineHtmlStrings)-i-1] = tmp
		}

		htmlStrings = append(htmlStrings, strings.Join(lineHtmlStrings, ""))
	}
	return strings.Join(htmlStrings, "")
}

func createMergedContent(curToken token.Token, parToken token.Token) string {
	content := []string{}
	switch parToken.ElmType {
	case token.LIST_ITEM:
		content = []string{"<li>", curToken.Content, "</li>"}
	case token.UL:
		content = []string{"<ul>", curToken.Content, "</ul>"}
	case token.OL:
		content = []string{"<ol>", curToken.Content, "</ol>"}
	case token.STRONG:
		content = []string{"<strong>", curToken.Content, "</strong>"}
	case token.H1:
		content = []string{"<h1>", curToken.Content, "</h1>"}
	case token.H2:
		content = []string{"<h2>", curToken.Content, "</h2>"}
	case token.H3:
		content = []string{"<h3>", curToken.Content, "</h3>"}
	case token.H4:
		content = []string{"<h4>", curToken.Content, "</h4>"}
	case token.H5:
		content = []string{"<h5>", curToken.Content, "</h5>"}
	case token.H6:
		content = []string{"<h6>", curToken.Content, "</h6>"}
	case token.MERGED:
		insertPos := getInsertPos(parToken.Content)
		content = []string{parToken.Content[0:insertPos], curToken.Content, parToken.Content[insertPos:]}
	}
	return strings.Join(content, "")
}

func getInsertPos(content string) int {
	// 親がmergedの時には必ずhtmlタグが存在するため、それを探索する
	htmlTagElmRegxp := HTML_TAG_REGEXP
	re := regexp.MustCompile(htmlTagElmRegxp)
	htmlTagIndexes := re.FindStringSubmatchIndex(content)
	insertPos := htmlTagIndexes[1] // htmlタグの終端の位置
	return insertPos
}
