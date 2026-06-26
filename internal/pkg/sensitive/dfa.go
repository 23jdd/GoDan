package sensitive

type DFAFilter struct {
	root *trieNode
}

type trieNode struct {
	children map[rune]*trieNode
	isEnd    bool
}

func NewDFAFilter(words []string) *DFAFilter {
	f := &DFAFilter{root: &trieNode{children: make(map[rune]*trieNode)}}
	for _, w := range words {
		f.AddWord(w)
	}
	return f
}

func (f *DFAFilter) AddWord(word string) {
	node := f.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &trieNode{children: make(map[rune]*trieNode)}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

func (f *DFAFilter) Check(text string) bool {
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		node := f.root
		for j := i; j < len(runes); j++ {
			child, ok := node.children[runes[j]]
			if !ok {
				break
			}
			node = child
			if node.isEnd {
				return true
			}
		}
	}
	return false
}

func (f *DFAFilter) Replace(text string, replace rune) string {
	runes := []rune(text)
	marked := make([]bool, len(runes))

	for i := 0; i < len(runes); i++ {
		node := f.root
		end := -1
		for j := i; j < len(runes); j++ {
			child, ok := node.children[runes[j]]
			if !ok {
				break
			}
			node = child
			if node.isEnd {
				end = j
			}
		}
		if end >= 0 {
			for k := i; k <= end; k++ {
				marked[k] = true
			}
		}
	}

	result := make([]rune, len(runes))
	for i, r := range runes {
		if marked[i] {
			result[i] = replace
		} else {
			result[i] = r
		}
	}
	return string(result)
}
