package token

type TokenType string

// 抽象構文木のノード
type Token struct {
	Id      int
	Parent  *Token
	ElmType TokenType
	Content string
}

// ノードの種類
const (
	TEXT      = "text"
	STRONG    = "strong"
	MERGED    = "merged"
	UL        = "ul"
	OL        = "ol"
	LIST_ITEM = "li"
	ROOT      = "root"
	H1        = "h1"
	H2        = "h2"
	H3        = "h3"
	H4        = "h4"
	H5        = "h5"
	H6        = "h6"
	A         = "a"
	A_HREF    = "a_href"
)

var HeaderTypeMap = map[int]TokenType{
	1: H1,
	2: H2,
	3: H3,
	4: H4,
	5: H5,
	6: H6,
}
