package internal

type tokenIterator struct {
	tokens []Token
	i      int
}

func (iter *tokenIterator) next() (Token, bool) {
	if iter.i == len(iter.tokens) {
		return Token{}, false
	}
	next := iter.tokens[iter.i]
	iter.i++
	return next, true
}

func (iter *tokenIterator) peek() (Token, bool) {
	if iter.i == len(iter.tokens) {
		return Token{}, false
	}
	return iter.tokens[iter.i], true
}

func (iter *tokenIterator) done() bool {
	return iter.i == len(iter.tokens)
}
