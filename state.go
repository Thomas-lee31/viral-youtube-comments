package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const pruneAge = 30 * 24 * time.Hour

type SeenVideos map[string]time.Time

func LoadState(path string) (SeenVideos, error) {
	seen := make(SeenVideos)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return seen, nil
		}
		return nil, fmt.Errorf("reading state file: %w", err)
	}

	if len(data) == 0 {
		return seen, nil
	}

	if err := json.Unmarshal(data, &seen); err != nil {
		return nil, fmt.Errorf("parsing state file: %w", err)
	}

	return seen, nil
}

// SaveState writes the state atomically: write to a temp file in the same
// directory, then rename over the target. This avoids corruption if the
// process is killed mid-write.
func SaveState(path string, seen SeenVideos) error {
	pruned := make(SeenVideos, len(seen))
	cutoff := time.Now().Add(-pruneAge)
	for id, t := range seen {
		if t.After(cutoff) {
			pruned[id] = t
		}
	}

	data, err := json.MarshalIndent(pruned, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling state: %w", err)
	}

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "seen_videos_*.tmp")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("renaming temp file: %w", err)
	}

	return nil
}
