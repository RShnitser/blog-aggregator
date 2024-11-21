package main

import(
	"blog-aggregator/internal/database"
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
)

func handleAddFeed(s *state, cmd command, user database.User) error{
	if len(cmd.args) < 2{
		return fmt.Errorf("feed name and url required")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
	})

	if err != nil{
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil{
		return err
	}

	fmt.Println(feed)
	return nil
}


func handleListFeeds(s *state, cmd command) error{
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil{
		return err
	}

	for _, feed := range feeds{
		fmt.Printf("Feed name: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Created By: %s\n", feed.UserName)
	}

	return nil
}