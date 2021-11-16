package splitter

import "strings"

// SimpleSplitter TODO
// PunctuationSplitter splitting keys by a specific set of symbols (i.e. punctuation)
type SimpleSplitter struct {
	divider rune
}

// NewSimpleSplitter  creates PunctuationSplitter
func NewSimpleSplitter(punctuation rune) *SimpleSplitter {
	return &SimpleSplitter{divider: punctuation}
}

// Split splits string key to fragments with given strategy
func (s *SimpleSplitter) Split(in string) []string {
	result := strings.Split(in, string(s.divider))
	for i, v := range result {
		// 包含数字, 非 pattern 字段, 替换掉
		if hasDigits(v) {
			result[i] = "*"
		}
	}
	return result
}

func hasDigits(s string) bool {
	b := false
	for _, c := range s {
		if c >= '0' && c <= '9' {
			b = true
			break
		}
	}
	return b
}
