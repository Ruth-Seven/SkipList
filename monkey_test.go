package skiplist

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"time"
)

// type operater interface {
// 	operate() error
// }

type operation func() error

type MonkeyTester struct {
	// Monkey []operater
	Monkey []operation
}

type Item struct {
	key   float64
	value int64
}

func (i Item) ExtendedKey() float64 {
	return i.key
}

type ListAndMap struct {
	list SkipList
	m    map[float64]int64
}

func NewListAndMap() ListAndMap {
	lm := ListAndMap{
		list: NewSkipList(0),
		m:    make(map[float64]int64),
	}
	return lm
}

func (lm *ListAndMap) genRandom() Item {
	value := randUint64(1000)
	key := float64((randUint64(10000))) // set the random distribution of key at (-10000, 10000)
	return Item{
		key:   key,
		value: int64(value),
	}
}

func (lm *ListAndMap) Insert() error {
	item := lm.genRandom()
	lm.list.Insert(item)
	lm.m[item.key] = item.value
	return nil
}

func (lm *ListAndMap) Delete() error {
	item := lm.genRandom()
	lm.list.Delete(item.ExtendedKey())
	delete(lm.m, item.key)
	return nil
}

func (lm *ListAndMap) Find() error {
	item := lm.genRandom()
	t1, ok1 := lm.list.Find(item.ExtendedKey())
	v, ok2 := lm.m[item.key]
	if ok1 != ok2 {
		return errors.New("not both empty")
	}
	if t1 != nil && v != t1.(Item).value {
		return errors.New("found value is not equal")
	}
	return nil
}

func (lm *ListAndMap) Equal() error {
	size := lm.list.Size()
	if size != len(lm.m) {
		return errors.New("the size is different")
	}
	for entity := lm.list.Minimal(); entity != nil; entity = lm.list.Next(entity) {
		if v, ok := lm.m[entity.key]; !ok {
			return errors.New("can't find item in map")
		} else if v != entity.Element.(Item).value {
			return errors.New("found value isn't identical")
		}
	}
	return nil
}

func NewMonkey() MonkeyTester {
	ml := NewListAndMap()
	monkey := MonkeyTester{
		Monkey: []operation{
			ml.Insert, ml.Delete, ml.Find, ml.Equal,
		},
	}
	return monkey
}

func (m *MonkeyTester) DoOnce() error {
	for _, op := range m.Monkey {
		err := op()
		if err != nil {
			return fmt.Errorf("operate A failed. op: %v, err: %v", op, err)
		}
	}
	return nil
}

func (m *MonkeyTester) RamdonDo() error {
	idx := randUint64(uint64(len(m.Monkey)))
	err := m.Monkey[idx]()
	if err != nil {
		return fmt.Errorf("fail to call RamdonDo . op: %T, err: %v", m.Monkey[idx], err)
	}
	return nil
}

func TestMain(t *testing.T) {
	testAmount := int64(1e6)
	gap := testAmount / 100
	threshold := int64(0)
	log.Print(testAmount)
	monkey := NewMonkey()
	err := monkey.DoOnce()
	if err != nil {
		t.Error(err)
		return
	}

	now := time.Now()
	for i := int64(0); i < testAmount; i++ {
		if i >= threshold {
			consume := time.Since(now).Seconds()
			log.Printf("Process:(%v/%v)%3f%% \t time: %8.3f， \t unit: %.6fs/op \n",
				i, testAmount, float64(i)/float64(testAmount)*100, consume, consume/float64(gap))
			threshold += gap
		}
		err = monkey.RamdonDo()
		if err != nil {
			t.Error(err)
			return
		}
	}
}
