package main

import (
	"strings"
	"testing"
)

func TestParseLLMResponse_FullResponse(t *testing.T) {
	input := `## Summary
This video covers Go concurrency patterns.
It explains goroutines and channels in depth.

## Timestamp Hero
- 0:00 Intro
- 2:30 Goroutines
- 5:00 Channels

## Inside Joke
"Just throw a goroutine at it" - famous last words.

## Value Add
The video skips over the scheduler internals, which are key to understanding why goroutines are cheap.`

	summary, comments := parseLLMResponse(input)

	if !strings.Contains(summary, "Go concurrency") {
		t.Errorf("summary missing expected content: %q", summary)
	}
	if len(comments) != 3 {
		t.Fatalf("expected 3 comments, got %d", len(comments))
	}
	if comments[0].name != "Timestamp Hero" {
		t.Errorf("first comment name = %q, want Timestamp Hero", comments[0].name)
	}
	if !strings.Contains(comments[0].value, "Goroutines") {
		t.Errorf("timestamp hero missing content: %q", comments[0].value)
	}
	if comments[1].name != "Inside Joke" {
		t.Errorf("second comment name = %q, want Inside Joke", comments[1].name)
	}
	if comments[2].name != "Value Add" {
		t.Errorf("third comment name = %q, want Value Add", comments[2].name)
	}
}

func TestParseLLMResponse_Empty(t *testing.T) {
	summary, comments := parseLLMResponse("")
	if summary != "" {
		t.Errorf("expected empty summary, got %q", summary)
	}
	if len(comments) != 0 {
		t.Errorf("expected no comments, got %d", len(comments))
	}
}

func TestParseLLMResponse_MissingSections(t *testing.T) {
	input := `## Summary
Just a summary, nothing else.`

	summary, comments := parseLLMResponse(input)
	if summary != "Just a summary, nothing else." {
		t.Errorf("unexpected summary: %q", summary)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(comments))
	}
}

func TestTruncateField_Short(t *testing.T) {
	got := truncateField("hello", 100)
	if got != "hello" {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestTruncateField_Exact(t *testing.T) {
	input := strings.Repeat("a", 1024)
	got := truncateField(input, 1024)
	if got != input {
		t.Errorf("expected unchanged at exact limit")
	}
}

func TestTruncateField_Over(t *testing.T) {
	input := strings.Repeat("a", 2000)
	got := truncateField(input, 1024)
	if len(got) != 1024 {
		t.Errorf("expected length 1024, got %d", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Error("expected ellipsis suffix")
	}
}
