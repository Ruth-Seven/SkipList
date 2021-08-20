package skiplist

type Element interface {
	ExtendedKey() float64
}

type EmptyElement int

func (e EmptyElement) ExtendedKey() float64 {
	return 0
}

var empty EmptyElement

type SkipItem struct {
	Element
	key     float64
	pre     *SkipItem
	forward []*SkipItem
	level   int
}
