package typetrie

import (
	"strings"

	"github.com/obukhov/redis-inventory/src/splitter"
	"github.com/obukhov/redis-inventory/src/trie"
)

// NewTypeTrie created Trie
func NewTypeTrie(splitter splitter.Splitter) *TypeTrie {
	node := trie.NewNode()
	node.AddAggregator(trie.NewAggregator())
	return &TypeTrie{
		root:     node,
		splitter: splitter,
	}
}

// TypeTrie stores data about keys in a prefix tree
type TypeTrie struct {
	root     *trie.Node
	splitter splitter.Splitter
}

// Add adds information about another key with set of params
func (t *TypeTrie) Add(key, keyType string, paramValues ...trie.ParamValue) {
	curNode := t.root
	var nextNode *trie.Node
	if childNode := curNode.GetChild(keyType); childNode == nil {
		nextNode = trie.NewNode()
		nextNode.AddAggregator(trie.NewAggregator())
		curNode.AddChild(keyType, nextNode)
	} else {
		nextNode = childNode
	}

	keyPieces := t.splitter.Split(key)
	pattern := strings.Join(keyPieces, t.splitter.Divider())
	var finalNode *trie.Node
	if childNode := nextNode.GetChild(pattern); childNode == nil {
		finalNode = trie.NewNode()
		finalNode.AddAggregator(trie.NewAggregator())
		nextNode.AddChild(pattern, finalNode)
	} else {
		finalNode = childNode
	}

	for _, p := range paramValues {
		curNode.Aggregator().Add(p.Param, p.Value)
		nextNode.Aggregator().Add(p.Param, p.Value)
		finalNode.Aggregator().Add(p.Param, p.Value)
	}
}

// Root returns root of the trie
func (t *TypeTrie) Root() *trie.Node {
	return t.root
}

// Clean TODO 清除 count ==1 的 pattern
func (t *TypeTrie) Clean() {
	for _, childNode := range t.root.Children {
		otherNode := trie.NewNode()
		otherNode.AddAggregator(trie.NewAggregator())
		childNode.AddChild("other", otherNode)
		for key, child := range childNode.Children {
			paramMap := child.Aggregator().Params
			if paramMap[trie.KeysCount] <= 1 {
				for k, v := range paramMap {
					otherNode.Aggregator().Add(k, v)
				}
				delete(childNode.Children, key)
			}
		}
	}
}
