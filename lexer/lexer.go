package lexer

import (
	"markdown-parser-go/token"
	"regexp"
	"strings"
)

func GenTextElement(id int, text string, parent token.Token) token.Token {
	return token.Token{Id: id, Parent: &parent, ElmType: token.TEXT, Content: text}
}

func GenStrongElement(id int, text string, parent token.Token) token.Token {
	return token.Token{Id: id, Parent: &parent, ElmType: token.STRONG, Content: ""}
}

func MatchWithStrongRegxp(text string) []int {
	re := regexp.MustCompile(token.STRONG_ELM_REGXP)
	return re.FindStringSubmatchIndex(text)
}

func MatchWithListRegxp(text string) []string {
	re := regexp.MustCompile(token.LIST_REGEXP)
	return re.FindStringSubmatch(text)
}

func Analize(markdown string) []string {
	state := token.NEUTRAL_STATE
	lists := []string{}
	mdArray := []string{}

	rawMdArray := regexp.MustCompile(`\r\n|\r|\n`).Split(markdown, -1)
	for index, md := range rawMdArray {
		// fmt.Println(md)
		listMatch := MatchWithListRegxp(md)
		if state == token.NEUTRAL_STATE && len(listMatch) != 0 {
			state = token.LIST_STATE
			lists = append(lists, strings.Join([]string{md, "\n"}, ""))
		} else if state == token.LIST_STATE && len(listMatch) != 0 {
			if index == len(rawMdArray)-1 {
				lists = append(lists, md)
				mdArray = append(mdArray, strings.Join(lists, ""))
			} else {
				lists = append(lists, strings.Join([]string{md, "\n"}, ""))
			}
		} else if state == token.LIST_STATE && len(listMatch) == 0 {
			state = token.NEUTRAL_STATE
			mdArray = append(mdArray, strings.Join(lists, ""))
			lists = []string{}
		}

		if len(lists) == 0 {
			mdArray = append(mdArray, md)
		}
	}
	return mdArray
}
