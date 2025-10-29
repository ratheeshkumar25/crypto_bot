package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type NewsArticle struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func FetchLatestNews(apiKey string) ([]NewsArticle, error) {
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?category=business&apiKey=%s", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Articles []NewsArticle `json:"articles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Articles, nil
}
