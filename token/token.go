package token

type Token struct {
	Id      int
	Parent  *Token
	ElmType string
	Content string
}

const (
	// tokenの状態
	TEXT   = "text"
	STRONG = "strong"
	MERGED = "merged"
	UL     = "ul"
	LIST   = "li"

	// 正規表現
	STRONG_ELM_REGXP = `\*\*(.*?)\*\*`
	LIST_REGEXP      = `(?m)^( *)([-|\*|\+] (.+))$`

	// リストの状態
	NEUTRAL_STATE = "neutral_state"
	LIST_STATE    = "list_state"
)
