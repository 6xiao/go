package TrieTree

type WordCount struct {
	Word  string
	Count int
}

type topN struct {
	top []*WordCount
	min int
}

func (this *topN) insert(r []rune, count int) {
	wc := &WordCount{string(r), count}
	for i := 0; i < len(this.top) && wc != nil; i++ {
		if this.top[i] != nil && this.top[i].Count == count && this.top[i].Word == wc.Word {
			return
		}
		if this.top[i] == nil || wc.Count > this.top[i].Count {
			this.top[i], wc = wc, this.top[i]
		}
	}
	if wc != nil {
		this.min = this.top[len(this.top)-1].Count
	}
}

func (this *topN) compact() []*WordCount {
	res := make([]*WordCount, 0, len(this.top))
	for _, v := range this.top {
		if v != nil {
			res = append(res, v)
		}
	}
	return res
}

type Node struct {
	Count    int
	Children map[rune]*Node
}

func NewTrieTree() *Node {
	return new(Node)
}

func (this *Node) add(seg []rune, index int, count int) int {
	if index >= len(seg) {
		this.Count += count
		return this.Count
	}

	if this.Children == nil {
		this.Children = make(map[rune]*Node, 1)
	}

	value := seg[index]
	if child, ok := this.Children[value]; !ok || child == nil {
		this.Children[value] = new(Node)
	}

	return this.Children[value].add(seg, index+1, count)
}

func (this *Node) Add(str string, count int) int {
	return this.add([]rune(str), 0, count)
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

func (this *Node) PrefixSearch(prefix string, topCount int) []*WordCount {
	seg, top := []rune(prefix), topN{make([]*WordCount, topCount), 0}
	if node := this.find(seg); node != nil {
		node.all(seg, &top)
	}
	return top.compact()
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

func (this *Node) SubstrSearch(sub string, topCount int) []*WordCount {
	seg, top := []rune(sub), topN{make([]*WordCount, topCount), 0}
	this.substr(this, nil, seg, &top)
	return top.compact()
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

func (this *Node) FuzzySearch(fuzzy string, topCount int) []*WordCount {
	seg, top := []rune(fuzzy), topN{make([]*WordCount, topCount), 0}
	this.fuzzy(nil, seg, 0, &top)
	return top.compact()
}
