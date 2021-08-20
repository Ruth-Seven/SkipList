![](https://github.com/actions/starter-workflows/blob/main/icons/go.svg)
# SkipList Implement in go
 
[![Go Report Card](https://goreportcard.com/badge/github.com/Ruth-Seven/skiplist)](https://goreportcard.com/report/github.com/Ruth-Seven/skiplist)   [![Go](https://github.com/Ruth-Seven/skiplist/actions/workflows/go.yaml/badge.svg)](https://github.com/Ruth-Seven/skiplist/actions/workflows/go.yaml)   [![Coverage Status](https://coveralls.io/repos/github/Ruth-Seven/skiplist/badge.svg)](https://coveralls.io/github/Ruth-Seven/skiplist)

A Basically implement of skiplist, which doesn't consider courrency in go.
skip lists： A probabilistaic alternative to balanced treees阐述了实现伪代码和设计上的工作

skip lists 的优点：

- 结构简单，每个节点占用的非信息空间少
- 查询速度和非递归自动调整树更快
- 每个节点的指针组创建之后就不必改动大小，只在插入或者删除后修改指向
- 有概率 P(Level=K)=1/(2^k)的概率，每个该节点的期望指针数是 2，修改和创建期望数也是如此

> 设立首次搜索 level= L(n)= log_2(n)，获取避免从当前最大 level 进行搜索，优化搜索次数。之后发现该条对性能没有明显提升，而且容易破坏更新值内容，故舍弃。

## drawback

limited by the lack of go generics, I have to use float64 to compare SkipItem. It's terrible, since you never know when you code will collision even if there is two rather different objet.

how to solve it? The answer is generics.

## test

```shell
 go clean -cache && go test -v
 go clean -cache && go test   -bench=. -v  # please don't use --race without runned parallel
 go clean -cache && go test   -bench=.  -v -cpuprofile=cpu.out -benchtime=100000000x
 go tool pprof -http localhost:8000 ./cpu.out
 go test -run=TestMonkey 
```

Results of benchmark:

```shell
#v0.0.2
// go clean -cache && go test  -run=NOTEST  -bench=.  -v -benchtime=1000000x -count=20 > x.out
// benchstat x.out
goos: darwin
goarch: amd64
pkg: github.com/MiniSkipList/skiplist
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

// -benchtime=1000000x  -count=20
SkipListInsert-12              141ns ± 3%
SkipListFind-12                143ns ± 6%
SkipListFindBiggerOrEqual-12   147ns ±13%
SkipListDelete-12             32.1ns ±20%
BucketmapInsert-12            71.9ns ±14%
BucketmapFind-12              66.5ns ±12%
BucketmapDelete-12            2.37ns ±18%

// -benchtime=10000000x -count=20
name                          time/op
SkipListInsert-12              177ns ± 4%
SkipListFind-12                174ns ± 7%
SkipListFindBiggerOrEqual-12   173ns ± 4%
SkipListDelete-12             26.3ns ± 3%
BucketmapInsert-12            88.0ns ± 2%
BucketmapFind-12              81.3ns ± 5%
BucketmapDelete-12            2.40ns ±36%

// -benchtime=100000000x -count=5
name                          time/op
SkipListInsert-12              227ns ±50%
SkipListFind-12                192ns ±10%
SkipListFindBiggerOrEqual-12   186ns ± 0%
SkipListDelete-12             25.5ns ± 1%
BucketmapInsert-12             122ns ± 2%
BucketmapFind-12               111ns ± 0%
BucketmapDelete-12            1.77ns ± 2%
```

operations / comparation times:

```shell
BenchmarkSkipListInsert findtime/ops: 4845312580 / 100000001 :48.453125
BenchmarkSkipListFind  findtime/ops: 5022697254 / 100000001 :50.226972
BenchmarkSkipListFindBiggerOrEqual findtime/ops: 5022697254 / 100000001 :50.226972
BenchmarkSkipListDelete findtime/ops: 2520589160 / 100000001 :25.205891
```

## Reference

[Skip Lists: A Probabilistic Alternative to Balanced Trees](https://15721.courses.cs.cmu.edu/spring2018/papers/08-oltpindexes1/pugh-skiplists-cacm1990.pdf)

[A Contention-Friendly, Non-Blocking Skip List](https://hal.inria.fr/hal-00699794v4/document)

[SkipList in C++](https://github.com/HiWong/SkipListPro)

[A another implement of skiplist in go](https://github.com/mauricegit/skiplist) And [open API](https://pkg.go.dev/github.com/MauriceGit/skiplist)
> The repo request you take a trunck of courage to read this codeeeeee.
