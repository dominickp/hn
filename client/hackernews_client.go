package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

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
	hackernewsURIPrefix = getEnvString("CHAN_HOST", defaultHackerNewsURIPrefix)
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
		Get(hackernewsURIPrefix + endpoint)
	if err != nil {
		return err
	}
	if response.IsError() {
		return errors.New(fmt.Sprintf("Error: %s", response.String()))
	}
	return nil
}

func GetTopStories() ([]int, error) {
	var topStories []int
	err := handleRequest("GET", "topstories.json", nil, &topStories)
	if err != nil {
		return nil, err
	}
	return topStories, nil
}

type Item struct {
	Id    int    `json:"id"`
	Type  string `json:"type"`
	By    string `json:"by"`
	Time  int    `json:"time"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
	Score int    `json:"score"`
	Kids  []int  `json:"kids"`
}

func GetItem(itemId int) (Item, error) {
	var item Item
	err := handleRequest("GET", fmt.Sprintf("item/%d.json", itemId), nil, &item)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

type TopMenuResponse struct {
	Items []Item `json:"items"`
}

func GetTopMenuResponse(maxItems int) (TopMenuResponse, error) {
	var topMenuResponse TopMenuResponse
	topStories, err := GetTopStories()
	if err != nil {
		return TopMenuResponse{}, err
	}
	for _, storyId := range topStories {
		item, err := GetItem(storyId)
		if err != nil {
			return TopMenuResponse{}, err
		}
		topMenuResponse.Items = append(topMenuResponse.Items, item)
		if len(topMenuResponse.Items) >= maxItems {
			break
		}
	}
	return topMenuResponse, nil

}
