package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

const systemPrompt = `You are a witty, insightful software engineer. Read the following YouTube video transcript. First, provide a 2-3 sentence summary of the video. Then, draft three distinct YouTube comments I could post:

1. The 'Timestamp Hero': A helpful bulleted summary of the key technical points with estimated timestamps.
2. The 'Inside Joke': A short, lighthearted, witty observation about a specific quote or moment.
3. The 'Value Add': An insightful, slightly contrarian take or deep-cut technical context that adds to the discussion.

Format your response exactly as:

## Summary
<summary here>

## Timestamp Hero
<comment here>

## Inside Joke
<comment here>

## Value Add
<comment here>`

const maxTranscriptWords = 12000

func GenerateComments(cfg *Config, video Video, transcript string) (string, error) {
	transcript = truncateWords(transcript, maxTranscriptWords)

	userMsg := fmt.Sprintf(
		"Video: %s\nChannel: %s\nURL: %s\n\nTranscript:\n%s",
		video.Title, video.ChannelName, video.URL, transcript,
	)

	client := openai.NewClient(
		option.WithAPIKey(cfg.OpenAIKey),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: shared.ChatModel(cfg.OpenAIModel),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userMsg),
		},
	})
	if err != nil {
		return "", fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("OpenAI returned no choices")
	}

	return completion.Choices[0].Message.Content, nil
}

func truncateWords(s string, max int) string {
	words := strings.Fields(s)
	if len(words) <= max {
		return s
	}
	return strings.Join(words[:max], " ") + " [transcript truncated]"
}
