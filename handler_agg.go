package main

import(
	"blog-aggregator/internal/database"
	"fmt"
	"time"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

func handleAggregate(s *state, cmd command) error{
	if len(cmd.args) < 1{
		return fmt.Errorf("time between requests required")
	}

	rawTime := cmd.args[0]
	timeBetweenRequests, err := time.ParseDuration(rawTime)
	if err != nil{
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func scrapeFeeds(s *state) error{
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil{
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ID:  feed.ID, LastFetchedAt: sql.NullTime{time.Now(), true}, UpdatedAt: time.Now()})

	currFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil{
		return err
	}
	for _, item := range currFeed.Channel.Item {
		//fmt.Printf("Found post: %s\n", item.Title)
		const layout = ""
		t, _ := time.Parse(layout, item.PubDate)
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: item.Title,
			Url: item.Link,
			Description: sql.NullString{item.Description, true},
			PublishedAt: t,
			FeedID: feed.ID,
		})
		if err != nil{
			fmt.Println(err)
		}
	}
	
	return nil
}