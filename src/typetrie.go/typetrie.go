package typetrie

import (
	"strings"

	"github.com/obukhov/redis-inventory/src/splitter"
	"github.com/obukhov/redis-inventory/src/trie"
)

// NewTypeTrie created Trie
func NewTypeTrie(splitter splitter.Splitter) *TypeTrie {
	return &TypeTrie{
		root:     make(map[string]KeyPatternMap),
		splitter: splitter,
	}
}

// TypeTrie TODO
// Trie stores data about keys in a prefix tree
type TypeTrie struct {
	root     map[string]KeyPatternMap
	splitter splitter.Splitter
}

// KeyPatternMap TODO
type KeyPatternMap struct {
	lengthMap map[int]SameLengthKeysSet
	Aggr      *trie.Aggregator
}

// SameLengthKeysSet TODO
type SameLengthKeysSet map[string]*trie.Aggregator

// // Value TODO
// type Value struct {
// 	// BytesSize size of the values in bytes
// 	BytesSize int64
// 	// KeysCount number of keys
// 	KeysCount int64
// }

// Add adds information about another key with set of params
func (t *TypeTrie) Add(key, keyType string, paramValues ...trie.ParamValue) {
	keyPieces := t.splitter.Split(key)
	pattern := strings.Join(keyPieces, t.splitter.Divider())
	length := len(keyPieces)
	var keyPatternMap KeyPatternMap
	var ok bool
	if keyPatternMap, ok = t.root[keyType]; !ok {
		t.root[keyType] = KeyPatternMap{
			lengthMap: make(map[int]SameLengthKeysSet),
			Aggr:      trie.NewAggregator(),
		}
		keyPatternMap = t.root[keyType]

	}
	typeAggregator := keyPatternMap.Aggr
	var sameLengthKeysSet SameLengthKeysSet
	if sameLengthKeysSet, ok = keyPatternMap.lengthMap[length]; !ok {
		keyPatternMap.lengthMap[length] = make(map[string]*trie.Aggregator)
		sameLengthKeysSet = keyPatternMap.lengthMap[length]
	}
	var aggregator *trie.Aggregator
	if aggregator, ok = sameLengthKeysSet[pattern]; !ok {
		sameLengthKeysSet[pattern] = trie.NewAggregator()
		aggregator = sameLengthKeysSet[pattern]
	}
	for _, p := range paramValues {
		aggregator.Add(p.Param, p.Value)
		typeAggregator.Add(p.Param, p.Value)
	}
}
