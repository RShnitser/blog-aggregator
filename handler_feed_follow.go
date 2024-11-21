package main

import(
	"blog-aggregator/internal/database"
	"fmt"
	"context"
	"github.com/google/uuid"
	"time"
)

func handleFollow(s *state, cmd command, user database.User) error{
	if len(cmd.args) < 1{
		return fmt.Errorf("url required")
	}

	feed, err := s.db.GetFeedFromUrl(context.Background(), cmd.args[0])
	if err != nil{
		return err
	}

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil{
		return err
	}

	fmt.Printf("Feed: %s\n", feed_follow.FeedName)	
	fmt.Printf("User: %s\n", feed_follow.UserName)	
	return nil
}

func handleListFollows(s *state, cmd command, user database.User) error{

	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil{
		return err
	}

	for _, feed_follow := range feed_follows{
		fmt.Printf("Feed: %s\n", feed_follow.FeedName)	
	}

	return nil

}

func handleUnfollow(s *state, cmd command, user database.User) error{

	if len(cmd.args) < 1{
		return fmt.Errorf("url required")
	}

	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{UserID: user.ID, Url:cmd.args[0]})
	if err != nil{
		return err
	}

	return nil
}