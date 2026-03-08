package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	rssFeedURL     = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"
	rssTimeout     = 10 * time.Second
	maxConcurrency = 10
)

// Atom feed structures for YouTube's RSS XML.
type AtomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Entries []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	VideoID   string     `xml:"videoId"`
	Title     string     `xml:"title"`
	Link      AtomLink   `xml:"link"`
	Author    AtomAuthor `xml:"author"`
	Published string     `xml:"published"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
}

type AtomAuthor struct {
	Name string `xml:"name"`
}

type Video struct {
	VideoID     string
	Title       string
	URL         string
	ChannelID   string
	ChannelName string
	Published   string
}

func FetchNewVideos(channels []Channel, seen SeenVideos) []Video {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		newVids []Video
		sem     = make(chan struct{}, maxConcurrency)
		client  = &http.Client{Timeout: rssTimeout}
	)

	for _, ch := range channels {
		wg.Add(1)
		sem <- struct{}{}

		go func(ch Channel) {
			defer wg.Done()
			defer func() { <-sem }()

			vid, err := fetchLatestVideo(client, ch)
			if err != nil {
				logf("RSS error for %s (%s): %v", ch.Name, ch.ChannelID, err)
				return
			}
			if vid == nil {
				return
			}

			mu.Lock()
			defer mu.Unlock()
			if _, already := seen[vid.VideoID]; !already {
				newVids = append(newVids, *vid)
			}
		}(ch)
	}

	wg.Wait()
	return newVids
}

func fetchLatestVideo(client *http.Client, ch Channel) (*Video, error) {
	url := fmt.Sprintf(rssFeedURL, ch.ChannelID)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var feed AtomFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("parsing XML: %w", err)
	}

	if len(feed.Entries) == 0 {
		return nil, nil
	}

	entry := feed.Entries[0]
	return &Video{
		VideoID:     entry.VideoID,
		Title:       entry.Title,
		URL:         entry.Link.Href,
		ChannelID:   ch.ChannelID,
		ChannelName: ch.Name,
		Published:   entry.Published,
	}, nil
}
