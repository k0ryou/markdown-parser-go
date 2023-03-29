package lexer

import (
	"markdown-parser-go/token"
	"regexp"
	"strings"
)

const (
	// 正規表現
	STRONG_ELM_REGXP = `\*\*(.*?)\*\*`
	UL_ITEM_REGXP    = `(?m)^( *)([-|\*|\+] (.+))$`
	OL_ITEM_REGXP    = `(?m)^( *)([0-9]+\. (.+))$`

	// markdownの現在の状態
	NEUTRAL_STATE = "neutral_state"
	UL_STATE      = "ul_state"
	OL_STATE      = "ol_state"
)

var textRegxpMap = map[token.TokenType]string{
	token.STRONG: STRONG_ELM_REGXP,
}

var listRegxpMap = map[token.TokenType]string{
	token.UL: UL_ITEM_REGXP,
	token.OL: OL_ITEM_REGXP,
}

func GenElementToken(id int, text string, parent token.Token, elmType token.TokenType) token.Token {
	return token.Token{Id: id, Parent: &parent, ElmType: elmType, Content: text}
}

func MatchIndexWithTextElmRegxp(text string, elmType token.TokenType) []int {
	re := regexp.MustCompile(textRegxpMap[elmType])
	return removeMinusVal(re.FindStringSubmatchIndex(text))
}

func MatchWithListElmRegxp(text string, elmType token.TokenType) []string {
	re := regexp.MustCompile(listRegxpMap[elmType])
	return re.FindStringSubmatch(text)
}

func removeMinusVal(slice []int) []int {
	res := []int{}
	for _, val := range slice {
		if val >= 0 {
			res = append(res, val)
		}
	}
	return res
}

func Analize(markdown string) []string {
	preState := NEUTRAL_STATE
	var nowState string
	lists := []string{}
	mdArray := []string{}

	rawMdArray := regexp.MustCompile(`\r\n|\r|\n`).Split(markdown, -1)
	for index, md := range rawMdArray {
		var isUlMatch bool = len(MatchWithListElmRegxp(md, token.UL)) > 0
		var isOlMatch bool = len(MatchWithListElmRegxp(md, token.OL)) > 0

		if isUlMatch {
			nowState = UL_STATE
		} else if isOlMatch {
			nowState = OL_STATE
		} else {
			nowState = NEUTRAL_STATE
		}

		if preState != nowState && len(lists) != 0 {
			appendLists2MdArray(&mdArray, &lists, md)
		}

		if nowState == UL_STATE || nowState == OL_STATE {
			appendListItem2Lists(&mdArray, &lists, md, index == len(rawMdArray)-1)
		} else {
			mdArray = append(mdArray, md)
		}
		preState = nowState
	}
	return mdArray
}

func appendListItem2Lists(mdArray *[]string, lists *[]string, md string, isLastLine bool) {
	*lists = append(*lists, strings.Join([]string{md, "\n"}, ""))
	if isLastLine {
		appendLists2MdArray(mdArray, lists, md)
	}
}

func appendLists2MdArray(mdArray *[]string, lists *[]string, md string) {
	*mdArray = append(*mdArray, strings.Join(*lists, ""))
	*lists = []string{}
}
