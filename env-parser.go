package debug

import (
	"regexp"
	"strings"
)

var splitRe = regexp.MustCompile(`\s+|,`)

type debugConfig struct {
	positive []string
	negative []string
}

func parseDebugEnv(envVar string) *debugConfig {
	result := new(debugConfig)
	names := splitRe.Split(envVar, -1)
	for _, name := range names {
		if name == "" {
			continue
		}

		if strings.HasPrefix(name, "-") {
			name = strings.TrimPrefix(name, "-")
			result.negative = append(result.negative, name)
		} else {
			result.positive = append(result.positive, name)
		}
	}
	return result
}

func (d *debugConfig) check(name string) bool {
	// Check negatives first
	for _, n := range d.negative {
		if compareWithMask(name, n) {
			return false
		}
	}

	for _, n := range d.positive {
		if compareWithMask(name, n) {
			return true
		}
	}

	return false
}

func compareWithMask(subj string, template string) bool {
	if template == "*" {
		return true
	}
	if strings.HasSuffix(template, "*") {
		prefix := strings.TrimSuffix(subj, "*")
		return strings.HasPrefix(subj, prefix)
	}
	return subj == template
}
