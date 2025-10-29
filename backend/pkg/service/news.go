package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type NewsArticle struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

var apiKey = `5768736115554f44968e3ff0449a4014`

func FetchLatestNews(apiKey string) ([]NewsArticle, error) {

	// 1. Using the /v2/everything endpoint for a more powerful search
	endpoint := "https://newsapi.org/v2/everything"

	keywords := "(crypto OR bitcoin OR ethereum) AND (market OR regulation OR price)"
	query := url.QueryEscape(keywords)

	// 3. Added 'sortBy=publishedAt' to get the newest articles first
	//    and 'language=en' to filter for English articles.
	fullURL := fmt.Sprintf("%s?q=%s&sortBy=publishedAt&language=en&apiKey=%s",
		endpoint,
		query,
		apiKey,
	)
	resp, err := http.Get(fullURL)
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
