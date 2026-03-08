package main

import (
	"strings"
	"testing"
)

func TestTruncateWords_Short(t *testing.T) {
	input := "one two three"
	got := truncateWords(input, 5)
	if got != input {
		t.Errorf("expected unchanged input, got %q", got)
	}
}

func TestTruncateWords_Exact(t *testing.T) {
	input := "one two three four five"
	got := truncateWords(input, 5)
	if got != input {
		t.Errorf("expected unchanged input, got %q", got)
	}
}

func TestTruncateWords_Truncated(t *testing.T) {
	input := "one two three four five six seven"
	got := truncateWords(input, 3)
	if got != "one two three [transcript truncated]" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestTruncateWords_Empty(t *testing.T) {
	got := truncateWords("", 10)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestTruncateWords_PreservesOriginalSpacing(t *testing.T) {
	input := "one  two\tthree"
	got := truncateWords(input, 100)
	if got != input {
		t.Errorf("expected original string preserved, got %q", got)
	}
}

func TestTruncateWords_LargeInput(t *testing.T) {
	words := make([]string, 15000)
	for i := range words {
		words[i] = "word"
	}
	input := strings.Join(words, " ")
	got := truncateWords(input, 12000)
	if !strings.HasSuffix(got, "[transcript truncated]") {
		t.Error("expected truncation suffix")
	}
	gotWords := strings.Fields(got)
	// 12000 words + 2 words from "[transcript truncated]"
	if len(gotWords) != 12002 {
		t.Errorf("expected 12002 words, got %d", len(gotWords))
	}
}
