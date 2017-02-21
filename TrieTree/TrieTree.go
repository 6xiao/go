package TrieTree

import (
	"sort"
)

type WordCount struct {
	Word  string
	Count int
}

type SearchResult []*WordCount

func (this SearchResult) Sort()              { sort.Sort(this) }
func (this SearchResult) Len() int           { return len(this) }
func (this SearchResult) Swap(i, j int)      { this[i], this[j] = this[j], this[i] }
func (this SearchResult) Less(i, j int) bool { return this[i].Count > this[j].Count }

type topN struct {
	res  SearchResult
	size int
	min  int
}

func newTopN(size int) *topN {
	if size <= 0 {
		size = 3
	}

	return &topN{make([]*WordCount, 0, size*4), size, 0}
}

func (this *topN) insert(r []rune, count int) {
	if count <= this.min {
		return
	}

	word := string(r)
	for _, wc := range this.res {
		if wc != nil && wc.Word == word {
			return
		}
	}

	this.res = append(this.res, &WordCount{word, count})
	if len(this.res) > this.size*3 {
		this.res.Sort()
		this.min = this.res[this.size].Count
		this.res = this.res[:this.size]
	}
}

func (this *topN) result() SearchResult {
	if len(this.res) > this.size {
		this.res.Sort()
		return this.res[:this.size]
	}
	return this.res
}

type Node struct {
	Count    int
	Children map[rune]*Node
}

func NewTrieTree() *Node {
	return new(Node)
}

func (this *Node) add(seg []rune, index, count, incr int) int {
	if index >= len(seg) {
		if count >= 0 {
			this.Count = count
		} else {
			this.Count += incr
		}
		return this.Count
	}

	if this.Children == nil {
		this.Children = make(map[rune]*Node, 1)
	}

	value := seg[index]
	if child, ok := this.Children[value]; !ok || child == nil {
		this.Children[value] = new(Node)
	}

	return this.Children[value].add(seg, index+1, count, incr)
}

func (this *Node) Add(str string, count, incr int) int {
	return this.add([]rune(str), 0, count, incr)
}

func (this *Node) all(seg []rune, top *topN) {
	if this.Count > top.min {
		top.insert(seg, this.Count)
	}

	for r, n := range this.Children {
		n.all(append(seg, r), top)
	}
}

func (this *Node) find(seg []rune) *Node {
	node := this
	for _, v := range seg {
		if child, ok := node.Children[v]; ok && child != nil {
			node = child
		} else {
			return nil
		}
	}
	return node
}

func (this *Node) PrefixSearch(prefix string, topCount int) SearchResult {
	seg, top := []rune(prefix), newTopN(topCount)
	if node := this.find(seg); node != nil {
		node.all(seg, top)
	}
	return top.result()
}

func (this *Node) substr(root *Node, pre, seg []rune, top *topN) {
	rp := append(pre, seg...)
	if node := root.find(rp); node != nil {
		node.all(rp, top)
	}

	for r, c := range this.Children {
		c.substr(root, append(pre, r), seg, top)
	}
}

func (this *Node) SubstrSearch(sub string, topCount int) SearchResult {
	seg, top := []rune(sub), newTopN(topCount)
	this.substr(this, nil, seg, top)
	return top.result()
}

func (this *Node) fuzzy(pre, seg []rune, index int, top *topN) {
	if index >= len(seg) {
		this.all(pre, top)
		for r, c := range this.Children {
			c.all(append(pre, r), top)
		}
		return
	}

	for r, c := range this.Children {
		if r == seg[index] {
			c.fuzzy(append(pre, r), seg, index+1, top)
		} else {
			c.fuzzy(append(pre, r), seg, index, top)
		}
	}
}

func (this *Node) FuzzySearch(fuzzy string, topCount int) SearchResult {
	seg, top := []rune(fuzzy), newTopN(topCount)
	this.fuzzy(nil, seg, 0, top)
	return top.result()
}
