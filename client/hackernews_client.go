package client

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/dominickp/hn/logger"
	"github.com/go-resty/resty/v2"
)

const (
	defaultHackerNewsURIPrefix = "https://hacker-news.firebaseio.com/v0/"
)

var (
	restyClient         *resty.Client
	hackernewsURIPrefix string
)

// getEnvString returns the value of the environment variable named by the key,
// or fallback if the environment variable is not set.
func getEnvString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func init() {
	hackernewsURIPrefix = getEnvString("HN_HOST", defaultHackerNewsURIPrefix)
	restyClient = resty.New().
		SetJSONMarshaler(json.Marshal).
		SetJSONUnmarshaler(json.Unmarshal).
		SetTimeout(time.Duration(5) * time.Second) // Set timeout to 5 seconds
}

// handleRequest is a helper function that handles the request to the 4channel API and captures fanout metrics.
func handleRequest(method string, endpoint string, headers map[string]string, result interface{}) error {
	response, err := restyClient.R().
		SetHeaders(headers).
		SetResult(result).
		Execute(method, hackernewsURIPrefix+endpoint)

	log.Logger.Printf("Request to %s%s returned %d", hackernewsURIPrefix, endpoint, response.StatusCode())
	if err != nil {
		return err
	}
	if response.IsError() {
		return fmt.Errorf("error: %s", response.String())
	}
	return nil
}

func GetTopStories() ([]int, error) {
	log.Logger.Println("Getting top stories")
	var topStories []int
	err := handleRequest("GET", "topstories.json", nil, &topStories)
	if err != nil {
		return nil, err
	}
	return topStories, nil
}

type Item struct {
	Id       int    `json:"id"`
	Type     string `json:"type"`
	By       string `json:"by"`
	Time     int    `json:"time"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	Url      string `json:"url"`
	Score    int    `json:"score"`
	Kids     []int  `json:"kids"`
	Comments []Item `json:"comments"`
}

func GetItem(itemId int) (Item, error) {
	log.Logger.Printf("Getting item %d", itemId)
	var item Item
	err := handleRequest("GET", fmt.Sprintf("item/%d.json", itemId), nil, &item)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func GetItemWithComments(itemId, maxComments int) (Item, error) {
	log.Logger.Printf("Getting item with comments %d", itemId)
	var item Item
	err := handleRequest("GET", fmt.Sprintf("item/%d.json", itemId), nil, &item)
	if err != nil {
		return Item{}, err
	}

	// Gather details of the comments
	for _, commentId := range item.Kids {
		if len(item.Comments) >= maxComments {
			break
		}
		log.Logger.Printf("Getting comment %d", commentId)
		var comment Item
		err := handleRequest("GET", fmt.Sprintf("item/%d.json", commentId), nil, &comment)
		if err != nil {
			return Item{}, err
		}

		if comment.Text == "" {
			// Skip comments with no text
			continue
		}

		if strings.HasPrefix(comment.Text, "[") {
			// Remove comments that are [dupe] or [dead] or [flagged]
			continue
		}

		item.Comments = append(item.Comments, comment)
	}

	return item, nil
}

type TopMenuResponse struct {
	Items []Item `json:"items"`
}

// Returns the top menu response with the top stories as items with only their IDs
func GetTopMenuResponse() (TopMenuResponse, error) {
	var topMenuResponse TopMenuResponse
	topStories, err := GetTopStories()
	if err != nil {
		return TopMenuResponse{}, err
	}

	topMenuResponse.Items = make([]Item, 0)
	for _, storyId := range topStories {
		topMenuResponse.Items = append(topMenuResponse.Items, Item{Id: storyId})
	}
	return topMenuResponse, nil

}

func (t TopMenuResponse) EnrichItems(pageSize, page int) {
	pageStories := t.Items[pageSize*(page-1) : pageSize*page]
	for i, item := range pageStories {
		if item.Type == "" {
			item, err := GetItem(item.Id)
			if err != nil {
				log.Logger.Printf("Error getting item %d: %v", item.Id, err)
				continue
			}
			pageStories[i] = item
		}
	}
}
