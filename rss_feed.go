package main

import(
	"context"
	"net/http"
	"io"
	"encoding/xml"
	"html"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil{
		return nil, err
	}
	req.Header.Add("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil{
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	rss := &RSSFeed{}
	err = xml.Unmarshal(data, rss)
	if err != nil{
		return nil, err
	}

	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for _, item := range rss.Channel.Item{
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return rss, nil
}