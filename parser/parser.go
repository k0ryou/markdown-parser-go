package parser

import (
	"markdown-parser-go/lexer"
	"markdown-parser-go/token"
	"regexp"
	"strings"
)

var rootToken = token.Token{
	Id:      0,
	Parent:  &token.Token{},
	ElmType: token.ROOT,
	Content: "",
}

func Parse(markdownRow string) []token.Token {
	if len(lexer.MatchWithUlItemRegxp(markdownRow)) != 0 {
		return tokenizeList(markdownRow, token.UL)
	}
	if len(lexer.MatchWithOlItemRegxp(markdownRow)) != 0 {
		return tokenizeList(markdownRow, token.OL)
	}
	initId := rootToken.Id
	return tokenizeText(&initId, rootToken, markdownRow)
}

func tokenizeText(id *int, p token.Token, text string) []token.Token {
	resultElements := []token.Token{}

	processingText := text
	parent := p
	for len(processingText) != 0 {
		matchIndexes := lexer.MatchWithStrongRegxp(processingText)

		if len(matchIndexes) == 0 {
			*id++
			onlyText := lexer.GenTextElement(*id, processingText, parent)
			processingText = ""
			resultElements = append(resultElements, onlyText)
		} else {
			matchTextStartIdx, matchTextEndIdx, innerTextStartIdx, innerTextEndIdx := matchIndexes[0], matchIndexes[1], matchIndexes[2], matchIndexes[3]
			matchText, innerText := processingText[matchTextStartIdx:matchTextEndIdx], processingText[innerTextStartIdx:innerTextEndIdx]
			// 先頭の通常文字
			if 0 < matchTextStartIdx {
				text := processingText[0:matchTextStartIdx]
				*id++
				textElm := lexer.GenTextElement(*id, text, parent)
				resultElements = append(resultElements, textElm)
				processingText = strings.Replace(processingText, text, "", 1)
			}

			// 太字
			*id++
			elm := lexer.GenStrongElement(*id, "", parent)
			parent = elm
			resultElements = append(resultElements, elm)
			processingText = strings.Replace(processingText, matchText, "", 1)
			resultElements = append(resultElements, tokenizeText(id, parent, innerText)...)
			parent = p
		}
	}

	return resultElements
}

func tokenizeList(listString string, listType string) []token.Token {
	id := 1
	rootListToken := token.Token{
		Id:      id,
		Parent:  &rootToken,
		ElmType: listType,
		Content: "",
	}
	parent := rootListToken
	tokens := []token.Token{rootListToken}

	listArray := regexp.MustCompile(`\r\n|\r|\n`).Split(listString, -1)
	for _, list := range listArray {
		var match = []string{}
		if listType == token.UL {
			match = lexer.MatchWithUlItemRegxp(list)
		} else {
			match = lexer.MatchWithOlItemRegxp(list)
		}

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
		listText := tokenizeText(&id, listToken, listInnerText)
		id += len(listText)
		tokens = append(tokens, listText...)
	}

	return tokens
}
