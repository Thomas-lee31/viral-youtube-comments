package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	discordTimeout       = 15 * time.Second
	discordFieldMaxChars = 1024
	discordEmbedColor    = 0xFF0000 // YouTube red
)

type discordPayload struct {
	Embeds []discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Title       string         `json:"title"`
	URL         string         `json:"url"`
	Color       int            `json:"color"`
	Description string         `json:"description,omitempty"`
	Fields      []discordField `json:"fields,omitempty"`
	Footer      *discordFooter `json:"footer,omitempty"`
}

type discordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type discordFooter struct {
	Text string `json:"text"`
}

func SendToDiscord(webhookURL string, video Video, llmResponse string) error {
	summary, comments := parseLLMResponse(llmResponse)

	embed := discordEmbed{
		Title: video.Title,
		URL:   video.URL,
		Color: discordEmbedColor,
		Footer: &discordFooter{
			Text: fmt.Sprintf("Channel: %s", video.ChannelName),
		},
	}

	if summary != "" {
		embed.Description = truncateField(summary, 4096)
	}

	for _, c := range comments {
		embed.Fields = append(embed.Fields, discordField{
			Name:  c.name,
			Value: truncateField(c.value, discordFieldMaxChars),
		})
	}

	payload := discordPayload{Embeds: []discordEmbed{embed}}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling discord payload: %w", err)
	}

	client := &http.Client{Timeout: discordTimeout}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("posting to discord: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

type comment struct {
	name  string
	value string
}

func parseLLMResponse(response string) (summary string, comments []comment) {
	sections := map[string]*string{
		"## Summary": &summary,
	}
	commentNames := []struct {
		header string
		label  string
	}{
		{"## Timestamp Hero", "Timestamp Hero"},
		{"## Inside Joke", "Inside Joke"},
		{"## Value Add", "Value Add"},
	}

	lines := strings.Split(response, "\n")
	var currentTarget *string
	var currentComment *comment

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if target, ok := sections[trimmed]; ok {
			currentTarget = target
			currentComment = nil
			continue
		}

		matched := false
		for _, cn := range commentNames {
			if trimmed == cn.header {
				comments = append(comments, comment{name: cn.label})
				currentComment = &comments[len(comments)-1]
				currentTarget = nil
				matched = true
				break
			}
		}
		if matched {
			continue
		}

		if currentTarget != nil {
			if *currentTarget != "" {
				*currentTarget += "\n"
			}
			*currentTarget += line
		} else if currentComment != nil {
			if currentComment.value != "" {
				currentComment.value += "\n"
			}
			currentComment.value += line
		}
	}

	summary = strings.TrimSpace(summary)
	for i := range comments {
		comments[i].value = strings.TrimSpace(comments[i].value)
	}

	return summary, comments
}

func truncateField(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
