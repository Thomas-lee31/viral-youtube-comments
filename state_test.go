package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadState_MissingFile(t *testing.T) {
	seen, err := LoadState("/tmp/nonexistent_state_test.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(seen) != 0 {
		t.Errorf("expected empty map, got %d entries", len(seen))
	}
}

func TestLoadState_EmptyFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "empty.json")
	os.WriteFile(tmp, []byte(""), 0644)

	seen, err := LoadState(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(seen) != 0 {
		t.Errorf("expected empty map, got %d entries", len(seen))
	}
}

func TestSaveAndLoadState_RoundTrip(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "state.json")

	original := SeenVideos{
		"abc123": time.Now().Add(-1 * time.Hour),
		"def456": time.Now(),
	}

	if err := SaveState(tmp, original); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	loaded, err := LoadState(tmp)
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded))
	}
	for id := range original {
		if _, ok := loaded[id]; !ok {
			t.Errorf("missing video ID %s after round-trip", id)
		}
	}
}

func TestSaveState_PrunesOldEntries(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "prune.json")

	seen := SeenVideos{
		"recent": time.Now().Add(-1 * time.Hour),
		"old":    time.Now().Add(-60 * 24 * time.Hour), // 60 days ago
	}

	if err := SaveState(tmp, seen); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	loaded, err := LoadState(tmp)
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}

	if _, ok := loaded["recent"]; !ok {
		t.Error("expected 'recent' to be kept")
	}
	if _, ok := loaded["old"]; ok {
		t.Error("expected 'old' to be pruned")
	}
}

func TestSaveState_AtomicWrite(t *testing.T) {
	dir := t.TempDir()
	tmp := filepath.Join(dir, "atomic.json")

	if err := SaveState(tmp, SeenVideos{"v1": time.Now()}); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".tmp" {
			t.Errorf("temp file %s was not cleaned up", e.Name())
		}
	}
}
