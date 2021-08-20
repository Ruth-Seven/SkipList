package skiplist

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type item1 struct {
	Key   float64
	Value int
}

func (i item1) ExtendedKey() float64 {
	return i.Key
}

func TestNewSkipList(t *testing.T) {
	checkpoint := assert.New(t)
	list := NewSkipList(RecommendedEps)
	checkpoint.NotNil(list)

	items := genTestItems(20, 1000)
	t.Log(items)
	for _, item := range items {
		list.Insert(item)
	}
	t.Log(list.String())

	_, ok := list.Find(items[3].Key)
	t.Log(list.String())
	checkpoint.True(ok)

	var element Element = item1{
		Key:   items[0].Key - 2,
		Value: 123,
	}
	newelement, ok := list.FindBiggerOrEqual(element.ExtendedKey())
	checkpoint.NotNil(newelement)
	item := newelement.(item1)
	t.Log(list.String())
	t.Log(element, newelement)
	checkpoint.Less(math.Abs(item.Key-items[0].Key), float64(2))
	checkpoint.True(ok)

	for i := range items {
		list.Delete(items[i].ExtendedKey())
		_, ok := list.Find(items[i].ExtendedKey())
		t.Log(list.String())
		checkpoint.False(ok)
	}
}

func TestGenRandomLevel(t *testing.T) {
	list := NewSkipList(RecommendedEps)
	distribution := make(map[int]int)
	for i := 0; i < 1e5; i++ {
		v := list.genRandomLevel()
		distribution[v]++
	}
	keys := make([]int, len(distribution))
	i := 0
	for k := range distribution {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	for _, k := range keys {
		t.Logf("%v : %v\n", k, distribution[k])
	}

	t.Logf("Note: please check manually: this output number should reduce a half each time.")
}

var list = NewSkipList(RecommendedEps) // init for benchmark
func BenchmarkSkipListInsert(t *testing.B) {
	item := item1{}
	for i := 0; i <= t.N; i++ {
		item.Key++
		list.Insert(item)
	}
}

func BenchmarkSkipListFind(t *testing.B) {
	item := item1{
		Key: -1,
	}
	for i := 0; i <= t.N; i++ {
		item.Key++
		list.Find(item.Key)
	}
}

func BenchmarkSkipListFindBiggerOrEqual(t *testing.B) {
	item := item1{
		Key: -1,
	}
	for i := 0; i <= t.N; i++ {
		item.Key++
		list.FindBiggerOrEqual(item.Key)
	}
}

func BenchmarkSkipListDelete(t *testing.B) {
	item := item1{}
	for i := 0; i <= t.N; i++ {
		item.Key++
		list.Delete(item.Key)
	}
}

var bucketmap = make(map[float64]int) // init for benchmark
func BenchmarkBucketmapInsert(t *testing.B) {
	item := item1{}
	for i := 0; i <= t.N; i++ {
		item.Key++
		bucketmap[item.Key] = item.Value
	}
}

func BenchmarkBucketmapFind(t *testing.B) {
	item := item1{
		Key: -1,
	}
	pass := make(chan int)
	for i := 0; i <= t.N; i++ {
		item.Key++
		if _, ok := bucketmap[item.Key]; ok {
			select {
			case quitStatus := <-pass:
				fmt.Print(quitStatus)
			default:
			}
		}
	}
}

func BenchmarkBucketmapDelete(t *testing.B) {
	item := item1{}
	for i := 0; i <= t.N; i++ {
		item.Key++
		delete(bucketmap, item.Key)
	}
}

func genTestItems(num, keyUpLimit int) []item1 {
	rand.Seed(time.Now().UnixNano())
	items := make([]item1, num)
	for i := 0; i < num; i++ {
		items[i] = item1{
			Key:   float64(randUint64(uint64(keyUpLimit))),
			Value: int(randUint64(uint64(num * 2))),
		}
	}
	return items
}

// TODO ADD map profilling tests.
