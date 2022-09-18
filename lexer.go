package expression

const (
	path tokenType = iota
	number
	boolean
	str
	ifNotFound
	minus
	plus
	sum
	leftParen
	rightParen
	times
	product
	divide
	length
	less
	lessEqual
	more
	moreEqual
	andOp
	and
	orOp
	or
	comma
	unknown
)

type tokenType int

type token struct {
	tokenType
	value interface{}
}

func lex(expr string) []tokenType {
	return nil
}
