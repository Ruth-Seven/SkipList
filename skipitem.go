package skiplist

import "constraint"

type Less interface {
	constraint.Ordered
}

type Key Less
type Value [T any]

var empty = 0

type SkipItem struct {
	Element Less
	key     float64
	pre     *SkipItem
	forward []*SkipItem
	level   int
}
