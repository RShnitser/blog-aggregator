package main

import(
	"blog-aggregator/internal/database"
	"fmt"
	"time"
	"context"
	"database/sql"
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
	fmt.Println("%v", currFeed)
	
	return nil
}