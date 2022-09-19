package internal

type stringIterator struct {
	runes []rune
	i     int
}

func (iter *stringIterator) next() (rune, bool) {
	if iter.i == len(iter.runes) {
		return 0, false
	}
	next := iter.runes[iter.i]
	iter.i++
	return next, true
}

func (iter *stringIterator) peek() (rune, bool) {
	if iter.i == len(iter.runes) {
		return 0, false
	}
	return iter.runes[iter.i], true
}

func (iter *stringIterator) done() bool {
	return iter.i == len(iter.runes)
}
