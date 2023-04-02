package parser

import (
	"markdown-parser-go/lexer"
	"markdown-parser-go/token"
	"regexp"
	"strings"
)

const (
	EOL_REGXP   = `\r\n|\r|\n`
	BLANK_REGXP = `[\s]+`

	MAX_SHARP_LEN = 6
)

var rootToken = token.Token{
	Id:      0,
	Parent:  &token.Token{},
	ElmType: token.ROOT,
	Content: "",
}

func Parse(markdownRow string) []token.Token {
	if len(lexer.MatchWithListElmRegxp(markdownRow, token.UL)) != 0 {
		return tokenizeList(markdownRow, token.UL)
	}
	if len(lexer.MatchWithListElmRegxp(markdownRow, token.OL)) != 0 {
		return tokenizeList(markdownRow, token.OL)
	}
	initId := rootToken.Id
	if matchHeaderList := lexer.MatchWithHeaderElmRegxp(markdownRow); isHeader(matchHeaderList) {
		return tokenizeHeader(&initId, rootToken, markdownRow, matchHeaderList)
	}
	return tokenizeText(&initId, rootToken, markdownRow)
}

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
			// 先頭の通常文字
			if 0 < matchTextStartIdx {
				text := processingText[0:matchTextStartIdx]
				*id++
				textElm := lexer.GenElementToken(*id, text, parent, token.TEXT)
				resultElements = append(resultElements, textElm)
				processingText = strings.Replace(processingText, text, "", 1)
			}

			// 太字
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

			if 0 < matchAnchorStartIdx {
				text := processingText[0:matchAnchorStartIdx]
				*id++
				textElm := lexer.GenElementToken(*id, text, parent, token.TEXT)
				resultElements = append(resultElements, textElm)
				processingText = strings.Replace(processingText, text, "", 1)
			}

			*id++
			elm := lexer.GenElementToken(*id, "", parent, token.A)
			parent = elm
			resultElements = append(resultElements, elm)
			processingText = strings.Replace(processingText, matchAnchorText, "", 1)

			// append hrefText
			*id++
			hrefEml := lexer.GenElementToken(*id, hrefText, parent, token.A_HREF)
			resultElements = append(resultElements, hrefEml)
			processingText = strings.Replace(processingText, hrefText, "", 1)

			// append anchorInnerText
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

func tokenizeList(listString string, listType token.TokenType) []token.Token {
	id := 1
	rootUlToken := token.Token{
		Id:      id,
		Parent:  &rootToken,
		ElmType: listType,
		Content: "",
	}
	parent := rootUlToken
	tokens := []token.Token{rootUlToken}

	listArray := regexp.MustCompile(EOL_REGXP).Split(listString, -1)
	for _, list := range listArray {
		match := lexer.MatchWithListElmRegxp(list, listType)
		if len(match) == 0 {
			continue
		}
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
		if matchHeaderList := lexer.MatchWithHeaderElmRegxp(listInnerText); isHeader(matchHeaderList) {
			listText = tokenizeHeader(&id, listToken, listInnerText, matchHeaderList)
		} else {
			listText = tokenizeText(&id, listToken, listInnerText)
		}
		// TODO: ここの挙動確認 id加算する必要あるか
		id += len(listText)
		tokens = append(tokens, listText...)
	}

	return tokens
}

func tokenizeHeader(id *int, parent token.Token, listString string, matchHeaderList []string) []token.Token {
	matchHeader := matchHeaderList[0]
	headerInnerText := matchHeaderList[3]

	sharpLen := len(regexp.MustCompile(BLANK_REGXP).Split(matchHeader, -1)[0])
	headerToken := token.HeaderTypeMap[sharpLen]

	*id++
	rootHeaderToken := token.Token{
		Id:      *id,
		Parent:  &parent,
		ElmType: headerToken,
		Content: "",
	}
	tokens := []token.Token{rootHeaderToken}
	listText := tokenizeText(id, rootHeaderToken, headerInnerText)
	tokens = append(tokens, listText...)

	return tokens
}

func isHeader(matchHeaderList []string) bool {
	if len(matchHeaderList) == 0 {
		return false
	}

	matchHeader := matchHeaderList[0]

	sharpLen := len(regexp.MustCompile(BLANK_REGXP).Split(matchHeader, -1)[0])

	return (0 < sharpLen && sharpLen <= MAX_SHARP_LEN)
}
