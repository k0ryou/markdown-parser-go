package generator

import (
	"markdown-parser-go/token"
	"math"
	"regexp"
	"sort"
	"strings"
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
			if curToken.Parent.ElmType != "root" {
				// 逆順なのでサイズから引く
				parIdx := parIdxMap[curToken.Parent.Id]
				parToken := rearrangedAst[parIdx]
				mergedToken := token.Token{
					Id:      parToken.Id,
					Parent:  parToken.Parent,
					ElmType: "merged",
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
	case "li":
		content = []string{"<li>", curToken.Content, "</li>"}
	case "ul":
		content = []string{"<ul>", curToken.Content, "</ul>"}
	case "strong":
		content = []string{"<strong>", curToken.Content, "</strong>"}
	case "merged":
		insertPos := getInsertPos(parToken.Content)
		content = []string{parToken.Content[0:insertPos], curToken.Content, parToken.Content[insertPos:]}
	}
	return strings.Join(content, "")
}

func getInsertPos(content string) int {
	// 親がmergedの時には必ずhtmlタグが存在するため、それを探索する
	htmlTagElmRegxp := `<(.*?)>`
	re := regexp.MustCompile(htmlTagElmRegxp)
	htmlTagIndexes := re.FindStringSubmatchIndex(content)
	insertPos := htmlTagIndexes[1] // htmlタグの終端の位置
	return insertPos
}
