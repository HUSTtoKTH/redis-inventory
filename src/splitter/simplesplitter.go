package splitter

import "strings"

// SimpleSplitter TODO
// PunctuationSplitter splitting keys by a specific set of symbols (i.e. punctuation)
type SimpleSplitter struct {
	divider string
}

// NewSimpleSplitter  creates PunctuationSplitter
func NewSimpleSplitter(punctuation string) *SimpleSplitter {
	return &SimpleSplitter{divider: punctuation}
}

// Split splits string key to fragments with given strategy
func (s *SimpleSplitter) Split(in string) []string {
	result := strings.Split(in, s.divider)
	for i, v := range result {
		// 包含数字, 非 pattern 字段, 替换掉
		if hasDigits(v) {
			result[i] = "*"
		}
	}
	return result
}

// Divider TODO
func (s *SimpleSplitter) Divider() string {
	return s.divider
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
