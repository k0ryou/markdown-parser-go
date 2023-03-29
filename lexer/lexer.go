package lexer

import (
	"markdown-parser-go/token"
	"regexp"
	"strings"
)

const (
	// 正規表現
	STRONG_ELM_REGXP = `\*\*(.*?)\*\*|__(.*?)__`
	UL_ITEM_REGXP    = `(?m)^( *)([-|\*|\+] (.+))$`
	OL_ITEM_REGXP    = `(?m)^( *)([0-9]+\. (.+))$`

	// markdownの現在の状態
	NEUTRAL_STATE = "neutral_state"
	UL_STATE      = "ul_state"
	OL_STATE      = "ol_state"
)

func GenTextElement(id int, text string, parent token.Token) token.Token {
	return token.Token{Id: id, Parent: &parent, ElmType: token.TEXT, Content: text}
}

func GenStrongElement(id int, text string, parent token.Token) token.Token {
	return token.Token{Id: id, Parent: &parent, ElmType: token.STRONG, Content: ""}
}

func MatchWithStrongRegxp(text string) []int {
	re := regexp.MustCompile(STRONG_ELM_REGXP)
	return removeMinusVal(re.FindStringSubmatchIndex(text))
}

func MatchWithUlItemRegxp(text string) []string {
	re := regexp.MustCompile(UL_ITEM_REGXP)
	return re.FindStringSubmatch(text)
}

func MatchWithOlItemRegxp(text string) []string {
	re := regexp.MustCompile(OL_ITEM_REGXP)
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
		var isUlMatch bool = len(MatchWithUlItemRegxp(md)) > 0
		var isOlMatch bool = len(MatchWithOlItemRegxp(md)) > 0

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
