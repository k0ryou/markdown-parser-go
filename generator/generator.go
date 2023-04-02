package generator

import (
	"markdown-parser-go/token"
	"math"
	"regexp"
	"sort"
	"strings"
)

const (
	// 正規表現
	HTML_TAG_REGEXP = `<(.*?)>`
	HREF_REGEXP     = `<a href='`
)

/*
抽象構文木からHTML形式のテキストを生成する
処理内容：
抽象構文木は行きかけ順にidが付与されているため、idの降順にノードを見ていき親ノードにマージしていく。
親ノードがルートの場合には、結果を保持する配列に格納する。
*/
func Generate(asts [][]token.Token) string {
	// 最終的な変換結果を持つ
	htmlStrings := []string{}
	// 抽象構文木を一つずつ取り出す
	for _, ast := range asts {
		// 行ごとの変換結果を持つ配列
		lineHtmlStrings := []string{}
		rearrangedAst := ast
		parIdxMap := map[int]int{}
		// 抽象構文木のノードをidの降順にソート
		sort.Slice(rearrangedAst, func(i, j int) bool { return rearrangedAst[i].Id > rearrangedAst[j].Id })
		// 各ノードについて、親ノードのidをメモしておく
		for idx, curToken := range rearrangedAst {
			parIdxMap[curToken.Id] = idx
		}
		for _, curToken := range rearrangedAst {
			if curToken.Parent.ElmType != token.ROOT { // 親ノードが根ではない場合
				parIdx := parIdxMap[curToken.Parent.Id]
				parToken := rearrangedAst[parIdx]
				// 親ノードと現在見ているノードをマージ
				mergedToken := token.Token{
					Id:      parToken.Id,
					Parent:  parToken.Parent,
					ElmType: token.MERGED,
					Content: createMergedContent(curToken, parToken),
				}
				// 親ノードにマージしたノードを格納
				rearrangedAst[parIdx] = mergedToken
			} else { // 親ノードが根である場合
				// 行ごとの結果に格納
				lineHtmlStrings = append(lineHtmlStrings, curToken.Content)
			}
		}

		// 行ごとの変換結果が逆順に並んでいるため反転する
		loopLimit := int(math.Sqrt(float64(len(lineHtmlStrings))))
		for i := 0; i < loopLimit; i++ {
			tmp := lineHtmlStrings[i]
			lineHtmlStrings[i] = lineHtmlStrings[len(lineHtmlStrings)-i-1]
			lineHtmlStrings[len(lineHtmlStrings)-i-1] = tmp
		}

		// 最終的な変換結果に追加
		htmlStrings = append(htmlStrings, strings.Join(lineHtmlStrings, ""))
	}
	return strings.Join(htmlStrings, "")
}

// curTokenとparTokenを結合したノードを生成する
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
	case token.A:
		content = []string{"<a href=''>", curToken.Content, "</a>"}
	case token.MERGED:
		// 親ノードが既にマージされていた場合、子ノードのテキストを追加する位置を探索する
		insertPos := getInsertPos(parToken.Content, curToken.ElmType)
		content = []string{parToken.Content[0:insertPos], curToken.Content, parToken.Content[insertPos:]}
	}
	return strings.Join(content, "")
}

// 親要素のcontent内から、子要素のテキストを追加する位置を探索する
func getInsertPos(content string, curElmType token.TokenType) int {
	var searchRegexp string
	if curElmType == token.A_HREF { // A_HREFの場合はaタグのhref='の位置を探索する
		searchRegexp = HREF_REGEXP
	} else { // それ以外の場合はhtmlタグ(<>)の終端の位置を探索する
		searchRegexp = HTML_TAG_REGEXP
	}
	re := regexp.MustCompile(searchRegexp)
	htmlTagIndexes := re.FindStringSubmatchIndex(content)
	insertPos := htmlTagIndexes[1]

	return insertPos
}
