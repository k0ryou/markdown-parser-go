package lexer

import (
	"markdown-parser-go/token"
	"regexp"
	"strings"
)

const (
	// 正規表現
	STRONG_ELM_REGXP  = `\*\*(.*)\*\*|__(.*)__`
	UL_ITEM_REGXP     = `(?m)^( *)([-|\*|\+] (.+))$`
	OL_ITEM_REGXP     = `(?m)^( *)([0-9]+\. (.+))$`
	HEADER_REGXP      = `(?m)^( *)(#+ (.+))$`
	PUNCTUATION_REGXP = `\pP`
	EOL_REGXP         = `[\r\n|\r|\n]+` // 複数行の空白に対応
	ANCHOR_REGXP      = `\[(.+)\]\((.+)\)`

	// markdownの現在の状態
	NEUTRAL_STATE = "neutral_state"
	UL_STATE      = "ul_state"
	OL_STATE      = "ol_state"
)

var listRegxpMap = map[token.TokenType]string{
	token.UL: UL_ITEM_REGXP,
	token.OL: OL_ITEM_REGXP,
}

func GenElementToken(id int, text string, parent token.Token, elmType token.TokenType) token.Token {
	return token.Token{Id: id, Parent: &parent, ElmType: elmType, Content: text}
}

func MatchIndexWithStrongElmRegxp(text string) []int {
	re := regexp.MustCompile(STRONG_ELM_REGXP)
	matchIndexList := removeMinusVal(re.FindStringSubmatchIndex(text))
	if len(matchIndexList) == 0 {
		return matchIndexList
	}

	matchTextStartIdx, matchTextEndIdx, innerTextStartIdx, innerTextEndIdx := matchIndexList[0], matchIndexList[1], matchIndexList[2], matchIndexList[3]
	var leftFrontChr byte = ' '
	leftBackChr := text[innerTextStartIdx]
	if 1 <= matchTextStartIdx {
		leftFrontChr = text[matchTextStartIdx-1]
	}

	var rightBackChr byte = ' '
	rightFrontChr := text[innerTextEndIdx-1]
	if matchTextEndIdx < len(text) {
		rightBackChr = text[matchTextEndIdx]
	}

	if !validLeftFlanking(leftFrontChr, leftBackChr) || !validRightFlanking(rightFrontChr, rightBackChr) { // **と__共通
		return []int{}
	}

	if text[matchTextStartIdx:innerTextStartIdx] == "__" { // __の場合
		// open tag
		if !(!validRightFlanking(leftFrontChr, leftBackChr) || (matchPunctuationChr(string(leftFrontChr)) && (leftBackChr == ' ' || matchPunctuationChr(string(leftBackChr))))) {
			return []int{}
		}

		// close tag
		if !(!validLeftFlanking(rightFrontChr, rightBackChr) || (matchPunctuationChr(string(rightBackChr)) && (rightFrontChr == ' ' || matchPunctuationChr(string(rightFrontChr))))) {
			return []int{}
		}
	}

	return matchIndexList
}

func MatchIndexWithAnchorElmRegxp(text string) []int {
	re := regexp.MustCompile(ANCHOR_REGXP)
	matchIndexList := removeMinusVal(re.FindStringSubmatchIndex(text))
	return matchIndexList
}

func MatchWithListElmRegxp(text string, elmType token.TokenType) []string {
	re := regexp.MustCompile(listRegxpMap[elmType])
	return re.FindStringSubmatch(text)
}

func MatchWithHeaderElmRegxp(text string) []string {
	re := regexp.MustCompile(HEADER_REGXP)
	return re.FindStringSubmatch(text)
}

func Analize(markdown string) []string {
	preState := NEUTRAL_STATE
	var nowState string
	lists := []string{}
	mdArray := []string{}

	rawMdArray := regexp.MustCompile(EOL_REGXP).Split(markdown, -1)
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

func validLeftFlanking(leftFrontChr byte, leftBackChr byte) bool {
	if leftBackChr == ' ' {
		return false
	}

	if !(!matchPunctuationChr(string(leftBackChr)) || (matchPunctuationChr(string(leftBackChr)) && (leftFrontChr == ' ' || matchPunctuationChr(string(leftFrontChr))))) {
		return false
	}

	return true
}

func validRightFlanking(rightFrontChr byte, rightBackChr byte) bool {
	if rightFrontChr == ' ' {
		return false
	}

	if !(!matchPunctuationChr(string(rightFrontChr)) || (matchPunctuationChr(string(rightFrontChr)) && (rightBackChr == ' ' || matchPunctuationChr(string(rightBackChr))))) {
		return false
	}

	return true
}

func matchPunctuationChr(chr string) bool {
	re := regexp.MustCompile(PUNCTUATION_REGXP)
	return re.MatchString(chr)
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
