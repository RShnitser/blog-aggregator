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
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil{
		return err
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	for _, feedFollow := range feedFollows{
		fmt.Printf("Feed: %s\n", feedFollow.FeedName)	
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