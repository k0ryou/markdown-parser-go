package lexer

import (
	"markdown-parser-go/token"
	"regexp"
	"strings"
)

const (
	// 正規表現
	STRONG_ELM_REGXP = `\*\*(.*?)\*\*|__(.*?)__`
	LIST_REGEXP      = `(?m)^( *)([-|\*|\+] (.+))$`

	// markdownの現在の状態
	NEUTRAL_STATE = "neutral_state"
	LIST_STATE    = "list_state"
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

func MatchWithListRegxp(text string) []string {
	re := regexp.MustCompile(LIST_REGEXP)
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
	state := NEUTRAL_STATE
	lists := []string{}
	mdArray := []string{}

	rawMdArray := regexp.MustCompile(`\r\n|\r|\n`).Split(markdown, -1)
	for index, md := range rawMdArray {
		// fmt.Println(md)
		listMatch := MatchWithListRegxp(md)
		if state == NEUTRAL_STATE && len(listMatch) != 0 {
			state = LIST_STATE
			lists = append(lists, strings.Join([]string{md, "\n"}, ""))
		} else if state == LIST_STATE && len(listMatch) != 0 {
			if index == len(rawMdArray)-1 {
				lists = append(lists, md)
				mdArray = append(mdArray, strings.Join(lists, ""))
			} else {
				lists = append(lists, strings.Join([]string{md, "\n"}, ""))
			}
		} else if state == LIST_STATE && len(listMatch) == 0 {
			state = NEUTRAL_STATE
			mdArray = append(mdArray, strings.Join(lists, ""))
			lists = []string{}
		}

		if len(lists) == 0 {
			mdArray = append(mdArray, md)
		}
	}
	return mdArray
}
