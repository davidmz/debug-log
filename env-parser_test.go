package debug

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareWithMask(t *testing.T) {
	testData := []struct {
		Subj     string
		Template string
		Result   bool
	}{
		{"aaa", "aaa", true},
		{"aaa", "aab", false},
		{"", "*", true},
		{"aaa", "*", true},
		{"aba", "*", true},
		{"ab", "ab*", true},
		{"aba", "ab*", true},
		{"abb", "ab*", true},
		{"abbc", "ab*", true},
	}
	for _, rec := range testData {
		t.Run(
			fmt.Sprintf("%q against %q should be %v", rec.Subj, rec.Template, rec.Result),
			func(t *testing.T) { assert.Equal(t, rec.Result, compareWithMask(rec.Subj, rec.Template)) },
		)
	}
}

func TestParseDebugEnv(t *testing.T) {
	testData := []struct {
		Env    string
		Result *debugConfig
	}{
		{"", &debugConfig{nil, nil}},
		{"  ", &debugConfig{nil, nil}},
		{"aaa", &debugConfig{[]string{"aaa"}, nil}},
		{"aaa,bbb", &debugConfig{[]string{"aaa", "bbb"}, nil}},
		{"aaa  bbb", &debugConfig{[]string{"aaa", "bbb"}, nil}},
		{"aaa, bbb", &debugConfig{[]string{"aaa", "bbb"}, nil}},
		{"aaa,-bbb", &debugConfig{[]string{"aaa"}, []string{"bbb"}}},
		{"aaa -bbb", &debugConfig{[]string{"aaa"}, []string{"bbb"}}},
		{"aaa,bbb,-bbb", &debugConfig{[]string{"aaa", "bbb"}, []string{"bbb"}}},
	}

	for _, rec := range testData {
		t.Run(
			fmt.Sprintf("%q should return %v", rec.Env, rec.Result),
			func(t *testing.T) { assert.Equal(t, rec.Result, parseDebugEnv(rec.Env)) },
		)
	}
}

func TestCheck(t *testing.T) {
	testData := []struct {
		Env    string
		Subj   string
		Result bool
	}{
		{"", "", false},
		{"aaa", "aaa", true},
		{"aaa", "aab", false},
		{"aaa,bbb", "aaa", true},
		{"aaa,bbb", "bbb", true},
		{"aaa,-bbb", "bbb", false},
		{"bb*,-bbb", "bbb", false},
		{"bb*,-bbb", "bba", true},
		{"*,-bbb", "bbb", false},
		{"*,-bbb", "aaa", true},
	}

	for _, rec := range testData {
		t.Run(
			fmt.Sprintf("%q against %q should return %v", rec.Subj, rec.Env, rec.Result),
			func(t *testing.T) { assert.Equal(t, rec.Result, parseDebugEnv(rec.Env).check(rec.Subj)) },
		)
	}
}
