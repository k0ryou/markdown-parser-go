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
	ROOT   = "root"
)
