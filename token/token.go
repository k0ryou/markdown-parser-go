package token

type TokenType string

type Token struct {
	Id      int
	Parent  *Token
	ElmType TokenType
	Content string
}

const (
	// tokenの状態
	TEXT      = "text"
	STRONG    = "strong"
	MERGED    = "merged"
	UL        = "ul"
	OL        = "ol"
	LIST_ITEM = "li"
	ROOT      = "root"
)
