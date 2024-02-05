// Intermediate layer between the Hacker News API and the application
package hn

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Story struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	User     string `json:"by"`
	Time     int    `json:"time"`
	Points   int    `json:"score"`
	Comments int    `json:"descendants"`
	Children []int  `json:"kids"`
}

type Comment struct {
	Text     string `json:"text"`
	User     string `json:"by"`
	Time     int    `json:"time"`
	Points   int    `json:"score"`
	Deleted  bool   `json:"deleted"`
	Children []int  `json:"kids"`
}

// FetchStory fetches a story from the Hacker News API
func FetchStory(id int) (story Story, err error) {
	// Create entrypoint
	entrypoint := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json?print=pretty", id)
	// Make request
	resp, err := http.Get(entrypoint)
	// Handle errors
	if err != nil {
		return story, err
	}
	if resp.StatusCode != 200 {
		return story, errors.New("status code is not 200")
	}
	// Close body after processing
	defer resp.Body.Close()
	// Unpack response
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return story, err
	}
	// Return
	return
}

// FetchStoryIds fetches a list of story ids from the Hacker News API
func FetchStoryIds(category, query string, frame [2]int) (ids []int, err error) {
	// Determine endpoint
	endpoint := map[string]string{
		"top":    "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty",
		"new":    "https://hacker-news.firebaseio.com/v0/newstories.json?print=pretty",
		"ask":    "https://hacker-news.firebaseio.com/v0/askstories.json?print=pretty",
		"show":   "https://hacker-news.firebaseio.com/v0/showstories.json?print=pretty",
		"jobs":   "https://hacker-news.firebaseio.com/v0/jobstories.json?print=pretty",
		"search": "http://hn.algolia.com/api/v1/search?query=%s&tags=story&hitsPerPage=500",
	}[category]
	// If search, extend with query
	if category == "search" {
		endpoint = fmt.Sprintf(endpoint, query)
	}
	// Make request
	resp, err := http.Get(endpoint)
	// Handle errors
	if err != nil {
		return ids, err
	}
	if resp.StatusCode != 200 {
		return ids, errors.New("status is not 200")
	}
	// Close body after processing
	defer resp.Body.Close()
	// Unpack story ids
	if category != "search" {
		err = json.NewDecoder(resp.Body).Decode(&ids)
		if err != nil {
			return ids, err
		}
	} else {
		var search struct {
			Hits []struct {
				ID string `json:"objectID"`
			} `json:"hits"`
		}
		err = json.NewDecoder(resp.Body).Decode(&search)
		if err != nil {
			return ids, err
		}
		for _, hit := range search.Hits {
			id, _ := strconv.Atoi(hit.ID)
			ids = append(ids, id)
		}
	}
	// Limit
	if len(ids) > frame[1] {
		ids = ids[frame[0]:frame[1]]
	}
	// Return
	return
}

// FetchComment fetches a comment from the Hacker News API
func FetchComment(id int) (comment Comment, err error) {
	// Create entrypoint
	entrypoint := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json?print=pretty", id)
	// Make request
	resp, err := http.Get(entrypoint)
	// Handle errors
	if err != nil {
		return comment, err
	}
	if resp.StatusCode != 200 {
		return comment, errors.New("status code is not 200")
	}
	// Close body after processing
	defer resp.Body.Close()
	// Unpack response
	if err := json.NewDecoder(resp.Body).Decode(&comment); err != nil {
		return comment, err
	}
	// Return
	return
}
