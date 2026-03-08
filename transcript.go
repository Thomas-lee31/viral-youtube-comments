package main

import (
	"fmt"
	"strings"

	"github.com/horiagug/youtube-transcript-api-go/pkg/yt_transcript"
)

func FetchTranscript(videoID string) (string, error) {
	client := yt_transcript.NewClient(
		yt_transcript.WithTimeout(30),
	)

	transcripts, err := client.GetTranscripts(videoID, []string{"en"})
	if err != nil {
		return "", fmt.Errorf("fetching transcript for %s: %w", videoID, err)
	}

	if len(transcripts) == 0 {
		return "", fmt.Errorf("no transcript available for %s", videoID)
	}

	var parts []string
	for _, t := range transcripts {
		for _, line := range t.Lines {
			text := strings.TrimSpace(line.Text)
			if text != "" {
				parts = append(parts, text)
			}
		}
	}

	if len(parts) == 0 {
		return "", fmt.Errorf("transcript for %s is empty", videoID)
	}

	return strings.Join(parts, " "), nil
}
