package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Channel struct {
	ChannelID string `json:"channel_id"`
	Name      string `json:"name"`
}

type Config struct {
	OpenAIKey         string
	OpenAIModel       string
	DiscordWebhookURL string
	ChannelsFile      string
	StateFile         string
	Channels          []Channel
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		OpenAIKey:         os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:       os.Getenv("OPENAI_MODEL"),
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		ChannelsFile:      os.Getenv("CHANNELS_FILE"),
		StateFile:         os.Getenv("STATE_FILE"),
	}

	if cfg.OpenAIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}
	if cfg.DiscordWebhookURL == "" {
		return nil, fmt.Errorf("DISCORD_WEBHOOK_URL is required")
	}
	if cfg.OpenAIModel == "" {
		cfg.OpenAIModel = "gpt-4o-mini"
	}
	if cfg.ChannelsFile == "" {
		cfg.ChannelsFile = "channels.json"
	}
	if cfg.StateFile == "" {
		cfg.StateFile = "seen_videos.json"
	}

	channels, err := loadChannels(cfg.ChannelsFile)
	if err != nil {
		return nil, fmt.Errorf("loading channels: %w", err)
	}
	cfg.Channels = channels

	return cfg, nil
}

func loadChannels(path string) ([]Channel, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var channels []Channel
	if err := json.Unmarshal(data, &channels); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	if len(channels) == 0 {
		return nil, fmt.Errorf("%s contains no channels", path)
	}

	return channels, nil
}
