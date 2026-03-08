package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	seen, err := LoadState(cfg.StateFile)
	if err != nil {
		log.Fatalf("state: %v", err)
	}

	logf("Loaded %d channels, %d previously seen videos", len(cfg.Channels), len(seen))

	newVideos := FetchNewVideos(cfg.Channels, seen)
	if len(newVideos) == 0 {
		logf("No new videos found. Exiting. (%s)", time.Since(start).Round(time.Millisecond))
		return
	}

	logf("Found %d new video(s)", len(newVideos))

	for _, video := range newVideos {
		logf("Processing: %s (%s)", video.Title, video.VideoID)

		transcript, err := FetchTranscript(video.VideoID)
		if err != nil {
			logf("Transcript unavailable for %s: %v (proceeding without)", video.VideoID, err)
			transcript = fmt.Sprintf("[Transcript unavailable] Video: %s by %s", video.Title, video.ChannelName)
		}

		llmResponse, err := GenerateComments(cfg, video, transcript)
		if err != nil {
			logf("LLM failed for %s: %v", video.VideoID, err)
			continue
		}

		if err := SendToDiscord(cfg.DiscordWebhookURL, video, llmResponse); err != nil {
			logf("Discord delivery failed for %s: %v (will retry next run)", video.VideoID, err)
			continue
		}

		seen[video.VideoID] = time.Now()
		if err := SaveState(cfg.StateFile, seen); err != nil {
			logf("Failed to save state after %s: %v", video.VideoID, err)
		}

		logf("Successfully processed and delivered: %s", video.Title)
	}

	logf("Done. Processed %d video(s) in %s", len(newVideos), time.Since(start).Round(time.Millisecond))
}

func logf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[youtubeads] "+format+"\n", args...)
}
