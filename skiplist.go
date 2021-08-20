package skiplist

import (
	"bytes"
	"math"
	"math/bits"
	"strconv"
)

const (
	RecommendedEps = 1e-9
	MaxLevel       = 40
)

type SkipList struct {
	head      *SkipItem
	tail      *SkipItem
	itemNum   int
	maxLevel  int
	level     int
	updated   []*SkipItem
	levelSeed uint64
	eps       float64
}

func NewSkipList(eps float64) SkipList {
	list := SkipList{
		maxLevel:  MaxLevel,
		updated:   make([]*SkipItem, MaxLevel),
		levelSeed: randUint64(math.MaxUint64),
		eps:       eps,
	}
	head := list.newNode(empty, false)
	list.head = &head
	list.tail = &head
	return list
}

func (list *SkipList) newNode(element Element, random bool) SkipItem {
	level := list.maxLevel
	if random {
		level = list.genRandomLevel()
	}
	return SkipItem{
		Element: element,
		key:     element.ExtendedKey(),
		forward: make([]*SkipItem, level+1),
		level:   level,
	}
}

func (list *SkipList) Find(key float64) (Element, bool) {
	node, ok := list.FindItem(key)
	if node == nil || !ok {
		return nil, false
	}
	return node.Element, true
}

func (list *SkipList) FindBiggerOrEqual(key float64) (Element, bool) {
	node, ok := list.FindBiggerOrEqualItem(key)
	if node == nil || !ok {
		return nil, false
	}
	return node.Element, true
}

func (list *SkipList) FindItem(key float64) (*SkipItem, bool) {
	node, _ := list.innerFind(key)
	if node == nil || !list.equal(key, node.key) {
		return nil, false
	}
	return node, true
}

func (list *SkipList) FindBiggerOrEqualItem(key float64) (*SkipItem, bool) {
	node, _ := list.innerFind(key)
	if node == nil || list.less(node.key, key) {
		return nil, false
	}
	return node, true
}

func (list *SkipList) Insert(element Element) {
	node, updatedLevel := list.innerFind(element.ExtendedKey())
	if node != nil && list.equal(element.ExtendedKey(), node.key) {
		node.Element = element
		return
	}

	newItem := list.newNode(element, true)
	newItem.pre = list.updated[0]
	for i := 0; i <= newItem.level; i++ {
		if updatedLevel < i {
			list.head.forward[i] = &newItem
		} else {
			newItem.forward[i] = list.updated[i].forward[i]
			list.updated[i].forward[i] = &newItem
		}
	}
	if newItem.forward[0] != nil {
		newItem.forward[0].pre = &newItem
	}
	if newItem.level > list.level {
		list.level = newItem.level
	}
	if newItem.forward[0] == nil {
		list.tail = &newItem
	}
	list.itemNum++
}

func (list *SkipList) Delete(key float64) bool {
	node, _ := list.innerFind(key)
	if node == nil {
		return true
	}
	if !list.equal(key, node.key) {
		return false
	}

	if node.forward[0] != nil {
		node.forward[0].pre = node.pre
	}
	node.pre = nil // Is it cleaning useless linker  helpful for gc?
	for i := node.level; i >= 0; i-- {
		if list.updated[i] == list.head && node.forward[i] == nil {
			list.level = i - 1
		}
		list.updated[i].forward[i] = node.forward[i]
	}
	if list.tail == node {
		list.tail = list.updated[0].forward[0]
	}
	list.itemNum--
	return true
}

// innerFind should find another element, which is equal as or biger than element
// the updatedLevel of returned list.UPDATED must bigger than or same as the level of selected item.
// that meaning
//
// Since we should keep correntness of the forword list, innnerFind gives up to compare last element
// of skiplist, which can speed up the time of looking up biggest element.
func (list *SkipList) innerFind(key float64) (node *SkipItem, updatedLevel int) {
	node = list.head
	updatedLevel = list.level
	for i := updatedLevel; i >= 0; i-- {
		for node.forward[i] != nil && node.forward[i].key < key {
			node = node.forward[i]
		}
		list.updated[i] = node
	}
	node = node.forward[0]
	return node, updatedLevel
}

// genRandomLevel
// @maxlevel: max level value must less than 64
func (list *SkipList) genRandomLevel() int {
	// tailzero := (uint64(1) << uint64(list.maxLevel)) | rand.Uint64() // TODO test the disturbtion
	// level := bits.TrailingZeros64(tailzero)
	// return level
	list.levelSeed++
	return bits.TrailingZeros64((uint64(1) << uint64(list.maxLevel)) | list.levelSeed)
}

func (list *SkipList) equal(a, b float64) bool {
	return !list.less(a, b) && !list.less(b, a)
	// return a == b
}

func (list *SkipList) less(a, b float64) bool {
	return a < b-list.eps
}

func (list *SkipList) String() string {
	buf := bytes.Buffer{}
	node := list.head
	buf.Write([]byte("Head"))
	for node != nil {
		node = node.forward[0]
		if node != nil {
			buf.Write([]byte(" -> " + strconv.FormatFloat(node.key, 'f', FormatLen, FormatBits)))
		}
	}
	buf.Write([]byte(" End"))
	return buf.String()
}

func (list *SkipList) Maximal() *SkipItem {
	if list.tail != nil {
		return list.tail
	}
	return nil
}

func (list *SkipList) Minimal() *SkipItem {
	if list.head.forward[0] != nil {
		return list.head.forward[0]
	}
	return nil
}

func (list *SkipList) Next(item *SkipItem) *SkipItem {
	if item == nil || item.forward[0] == nil {
		return nil
	}
	return item.forward[0]
}

func (list *SkipList) Before(item *SkipItem) *SkipItem {
	if item == nil || item.pre == nil || item.pre == list.head {
		return nil
	}
	return item.pre
}

func (list *SkipList) Empty() bool {
	return list.itemNum == 0
}

func (list *SkipList) Size() int {
	return list.itemNum
}
