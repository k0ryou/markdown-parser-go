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

	// Markdownの現在見ている行の状態
	NEUTRAL_STATE = "neutral_state"
	UL_STATE      = "ul_state"
	OL_STATE      = "ol_state"
)

var listRegxpMap = map[token.TokenType]string{
	token.UL: UL_ITEM_REGXP,
	token.OL: OL_ITEM_REGXP,
}

// 指定された内容のTokenを生成する
func GenElementToken(id int, text string, parent token.Token, elmType token.TokenType) token.Token {
	return token.Token{Id: id, Parent: &parent, ElmType: elmType, Content: text}
}

// text中の強調タグ(strong)の正規表現にマッチするすべての位置を取得する
func MatchIndexWithStrongElmRegxp(text string) []int {
	re := regexp.MustCompile(STRONG_ELM_REGXP)
	matchIndexList := removeMinusVal(re.FindStringSubmatchIndex(text))
	if len(matchIndexList) == 0 {
		return matchIndexList
	}

	matchTextStartIdx, matchTextEndIdx, innerTextStartIdx, innerTextEndIdx := matchIndexList[0], matchIndexList[1], matchIndexList[2], matchIndexList[3]
	// <strong>(open tag)の両端のbyteを取得
	var leftFrontChr byte = ' '
	leftBackChr := text[innerTextStartIdx]
	if 1 <= matchTextStartIdx {
		leftFrontChr = text[matchTextStartIdx-1]
	}

	// </strong>(close tag)の両端のbyteを取得
	var rightBackChr byte = ' '
	rightFrontChr := text[innerTextEndIdx-1]
	if matchTextEndIdx < len(text) {
		rightBackChr = text[matchTextEndIdx]
	}

	// 以下バリデーション(条件を満たさなかった場合は空配列を返す)
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

// text中のAnchorタグ(<a>)の正規表現にマッチするすべての位置取得する
func MatchIndexWithAnchorElmRegxp(text string) []int {
	re := regexp.MustCompile(ANCHOR_REGXP)
	return removeMinusVal(re.FindStringSubmatchIndex(text))
}

// text中のHeaderタグ(h1~h6)の正規表現にマッチするすべての文字列を取得する
func MatchWithHeaderElmRegxp(text string) []string {
	re := regexp.MustCompile(HEADER_REGXP)
	return re.FindStringSubmatch(text)
}

// text中のListタグ(ul, ol)の正規表現にマッチするすべての文字列を取得する
func MatchWithListElmRegxp(text string, elmType token.TokenType) []string {
	re := regexp.MustCompile(listRegxpMap[elmType])
	return re.FindStringSubmatch(text)
}

// markdownを適切な区間ごとに分割する(主にul,olを余分に作らないことが目的)
func Analize(markdown string) []string {
	// 前の状態と現在の状態を持つ
	preState := NEUTRAL_STATE
	var nowState string
	// 一行のデータを持つ配列
	lists := []string{}
	// 分割された結果を持つ配列
	mdArray := []string{}

	// 改行ごとに分割
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

		// 前の状態と異なる場合、listsをmdArrayに保存してリセット
		if preState != nowState && len(lists) != 0 {
			appendLists2MdArray(&mdArray, &lists, md)
		}

		// 現在の状態がリストの要素(li)ならlistsにデータを追加、それ以外ならmdArrayに直接追加
		if nowState == UL_STATE || nowState == OL_STATE {
			appendListItem2Lists(&mdArray, &lists, md, index == len(rawMdArray)-1)
		} else {
			mdArray = append(mdArray, md)
		}
		preState = nowState
	}
	return mdArray
}

// listsにリストの要素を追加
func appendListItem2Lists(mdArray *[]string, lists *[]string, md string, isLastLine bool) {
	*lists = append(*lists, strings.Join([]string{md, "\n"}, ""))
	if isLastLine { // 最後の行ならlistsを保存する
		appendLists2MdArray(mdArray, lists, md)
	}
}

// listsをmdArrayに追加する
func appendLists2MdArray(mdArray *[]string, lists *[]string, md string) {
	*mdArray = append(*mdArray, strings.Join(*lists, ""))
	*lists = []string{}
}

// タグの左側のバリデーション
func validLeftFlanking(leftFrontChr byte, leftBackChr byte) bool {
	if leftBackChr == ' ' {
		return false
	}

	if !(!matchPunctuationChr(string(leftBackChr)) || (matchPunctuationChr(string(leftBackChr)) && (leftFrontChr == ' ' || matchPunctuationChr(string(leftFrontChr))))) {
		return false
	}

	return true
}

// タグの右側のバリデーション
func validRightFlanking(rightFrontChr byte, rightBackChr byte) bool {
	if rightFrontChr == ' ' {
		return false
	}

	if !(!matchPunctuationChr(string(rightFrontChr)) || (matchPunctuationChr(string(rightFrontChr)) && (rightBackChr == ' ' || matchPunctuationChr(string(rightBackChr))))) {
		return false
	}

	return true
}

// chrが区切り文字であるか判別する
func matchPunctuationChr(chr string) bool {
	re := regexp.MustCompile(PUNCTUATION_REGXP)
	return re.MatchString(chr)
}

// sliceから負の値を削除する
func removeMinusVal(slice []int) []int {
	res := []int{}
	for _, val := range slice {
		if val >= 0 {
			res = append(res, val)
		}
	}
	return res
}
