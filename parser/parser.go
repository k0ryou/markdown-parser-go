package parser

import (
	"markdown-parser-go/lexer"
	"markdown-parser-go/token"
	"regexp"
	"strings"
)

const (
	// 正規表現
	EOL_REGXP   = `\r\n|\r|\n`
	BLANK_REGXP = `[\s]+`

	// headerの#の最大長
	MAX_HEADER_SHARP_LEN = 6
)

// 抽象構文木の根
var rootToken = token.Token{
	Id:      0,
	Parent:  &token.Token{},
	ElmType: token.ROOT,
	Content: "",
}

/*
markdownRowを抽象構文木に変換する
処理内容：
深さ優先探索により抽象構文木を生成する。
まず、rootTokenを根ノードとする。
次に、以下の優先順位でタグを探索する。
1. ul または ol
2. h1~h6
3. その他
タグが見つかったらノードを生成し、親ノードと紐づける。
ノードを生成した部分の文字列を除き、再帰的に同じ処理を繰り返す。
文字列が空になったら処理を終了する。
また、行きかけ順にidを付与する(HTMLへの変換時に用いる)。
*/
func Parse(markdownRow string) []token.Token {
	if len(lexer.MatchWithListElmRegxp(markdownRow, token.UL)) != 0 {
		return tokenizeList(markdownRow, token.UL)
	}
	if len(lexer.MatchWithListElmRegxp(markdownRow, token.OL)) != 0 {
		return tokenizeList(markdownRow, token.OL)
	}
	curId := rootToken.Id
	if matchHeaderList := lexer.MatchWithHeaderElmRegxp(markdownRow); isHeader(matchHeaderList) {
		return tokenizeHeader(&curId, rootToken, markdownRow, matchHeaderList)
	}
	return tokenizeText(&curId, rootToken, markdownRow)
}

// テキスト(strong, a, ...)を抽象構文木に変換する
func tokenizeText(id *int, p token.Token, text string) []token.Token {
	resultElements := []token.Token{}

	processingText := text
	parent := p
	for len(processingText) != 0 {
		matchStrongIndexes := lexer.MatchIndexWithStrongElmRegxp(processingText)
		matchAnchorIndexes := lexer.MatchIndexWithAnchorElmRegxp(processingText)

		if len(matchStrongIndexes) != 0 { // 太字の要素
			matchTextStartIdx, matchTextEndIdx, innerTextStartIdx, innerTextEndIdx := matchStrongIndexes[0], matchStrongIndexes[1], matchStrongIndexes[2], matchStrongIndexes[3]
			matchText, innerText := processingText[matchTextStartIdx:matchTextEndIdx], processingText[innerTextStartIdx:innerTextEndIdx]
			// 先頭の通常文字の処理
			if 0 < matchTextStartIdx {
				text := processingText[0:matchTextStartIdx]
				*id++
				textElm := lexer.GenElementToken(*id, text, parent, token.TEXT)
				resultElements = append(resultElements, textElm)
				processingText = strings.Replace(processingText, text, "", 1)
			}

			// 太字の処理
			*id++
			elm := lexer.GenElementToken(*id, "", parent, token.STRONG)
			parent = elm
			resultElements = append(resultElements, elm)
			processingText = strings.Replace(processingText, matchText, "", 1)
			resultElements = append(resultElements, tokenizeText(id, parent, innerText)...)
			parent = p
		} else if len(matchAnchorIndexes) != 0 { // aタグの要素
			matchAnchorStartIdx, matchAnchorEndIdx, innerTextStartIdx, innerTextEndIdx, hrefTextStartIdx, hrefTextEndIdx := matchAnchorIndexes[0], matchAnchorIndexes[1], matchAnchorIndexes[2], matchAnchorIndexes[3], matchAnchorIndexes[4], matchAnchorIndexes[5]
			// [anchorInnerText](hrefText)
			matchAnchorText, anchorInnerText, hrefText := processingText[matchAnchorStartIdx:matchAnchorEndIdx], processingText[innerTextStartIdx:innerTextEndIdx], processingText[hrefTextStartIdx:hrefTextEndIdx]
			// 先頭の通常文字
			if 0 < matchAnchorStartIdx {
				text := processingText[0:matchAnchorStartIdx]
				*id++
				textElm := lexer.GenElementToken(*id, text, parent, token.TEXT)
				resultElements = append(resultElements, textElm)
				processingText = strings.Replace(processingText, text, "", 1)
			}

			// aタグの処理
			*id++
			elm := lexer.GenElementToken(*id, "", parent, token.A)
			parent = elm
			resultElements = append(resultElements, elm)
			processingText = strings.Replace(processingText, matchAnchorText, "", 1)

			// aタグのリンクテキストの処理
			*id++
			hrefEml := lexer.GenElementToken(*id, hrefText, parent, token.A_HREF)
			resultElements = append(resultElements, hrefEml)
			processingText = strings.Replace(processingText, hrefText, "", 1)

			// aタグ内のテキストの処理
			resultElements = append(resultElements, tokenizeText(id, parent, anchorInnerText)...)
			parent = p
		} else {
			*id++
			onlyText := lexer.GenElementToken(*id, processingText, parent, token.TEXT)
			processingText = ""
			resultElements = append(resultElements, onlyText)
		}
	}

	return resultElements
}

// リストを抽象構文木に変換
func tokenizeList(listString string, listType token.TokenType) []token.Token {
	// リストの根を生成
	id := 1
	rootUlToken := token.Token{
		Id:      id,
		Parent:  &rootToken,
		ElmType: listType,
		Content: "",
	}
	parent := rootUlToken
	tokens := []token.Token{rootUlToken}

	// リストを改行ごとに分割する
	listArray := regexp.MustCompile(EOL_REGXP).Split(listString, -1)
	for _, list := range listArray {
		// リストの要素を取得する
		match := lexer.MatchWithListElmRegxp(list, listType)
		if len(match) == 0 {
			continue
		}

		// リストの要素を生成
		id++
		listToken := token.Token{
			Id:      id,
			Parent:  &parent,
			ElmType: token.LIST_ITEM,
			Content: "",
		}
		tokens = append(tokens, listToken)
		listInnerText := match[3]
		var listText []token.Token
		// リストの要素内にヘッダータグがあった場合は処理する
		if matchHeaderList := lexer.MatchWithHeaderElmRegxp(listInnerText); isHeader(matchHeaderList) {
			listText = tokenizeHeader(&id, listToken, listInnerText, matchHeaderList)
		} else {
			listText = tokenizeText(&id, listToken, listInnerText)
		}
		tokens = append(tokens, listText...)
	}

	return tokens
}

// ヘッダーを抽象構文木に変換する
func tokenizeHeader(id *int, parent token.Token, listString string, matchHeaderList []string) []token.Token {
	matchHeader := matchHeaderList[0]
	headerInnerText := matchHeaderList[3]

	sharpLen := len(regexp.MustCompile(BLANK_REGXP).Split(matchHeader, -1)[0])
	headerToken := token.HeaderTypeMap[sharpLen]

	// ヘッダーの根を生成
	*id++
	rootHeaderToken := token.Token{
		Id:      *id,
		Parent:  &parent,
		ElmType: headerToken,
		Content: "",
	}
	tokens := []token.Token{rootHeaderToken}
	// ヘッダー内のテキストについて処理
	listText := tokenizeText(id, rootHeaderToken, headerInnerText)
	tokens = append(tokens, listText...)

	return tokens
}

// ヘッダーの正規表現にマッチした文字列がヘッダータグであるか確認する(先頭の#の数を確認する)
func isHeader(matchHeaderList []string) bool {
	if len(matchHeaderList) == 0 {
		return false
	}

	matchHeader := matchHeaderList[0]

	sharpLen := len(regexp.MustCompile(BLANK_REGXP).Split(matchHeader, -1)[0])

	return (0 < sharpLen && sharpLen <= MAX_HEADER_SHARP_LEN)
}
