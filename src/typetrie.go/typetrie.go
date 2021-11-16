package typetrie

// NewTypeTrie TODO
// NewTrie created Trie
func NewTypeTrie(splitter Splitter, maxChildren int) *Trie {
	node := NewNode()
	node.AddAggregator(NewAggregator())

	return &Trie{
		root:        node,
		splitter:    splitter,
		maxChildren: maxChildren,
	}
}

// Trie stores data about keys in a prefix tree
type Trie struct {
	root        *Node
	splitter    Splitter
	maxChildren int
}
