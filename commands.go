package main

import(
	"blog-aggregator/internal/database"
	"blog-aggregator/internal/config"
	"fmt"
	"github.com/google/uuid"
	"time"
	"context"
)

type state struct{
	db  *database.Queries
	cfg *config.Config
}

type command struct{
	name string
	args []string
}

type commands struct{
	 cmdMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error){
	c.cmdMap[name] = f
}

func (c *commands) run(s *state, cmd command) error{
	f, ok := c.cmdMap[cmd.name]
	if ok{
		return f(s, cmd)
	}
	return fmt.Errorf("invalid command")
}

func handlerLogin(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("username required")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil{
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil{
		return err
	}

	fmt.Printf("user name set to %s\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("username required")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
	})
	if err != nil{
		return err
	}

	fmt.Printf("created user with name %s\n", user.Name)

	err = s.cfg.SetUser(user.Name)
	if err != nil{
		return err
	}

	return err
}

func handleResetDatabase(s *state, cmd command) error{
	err := s.db.DeleteUsers(context.Background())
	if err != nil{
		return err
	}
	fmt.Println("All users deleted")
	return nil
}

func handleListUsers(s *state, cmd command) error{
	users, err := s.db.GetUsers(context.Background())
	if err != nil{
		return err
	}

	current := s.cfg.CurrentUserName

	for _, user := range users{
		if user.Name == current{
			fmt.Printf("* %s (current)\n", user.Name)
		}else{
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}







func handleFollowFeed(s *state, cmd command, user database.User) error{
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

