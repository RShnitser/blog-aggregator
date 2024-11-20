package main

import(
	"blog-aggregator/internal/database"
	"blog-aggregator/internal/config"
	"fmt"
	"github.com/google/uuid"
	"time"
	"context"
	"net/http"
	"io"
	"encoding/xml"
	"html"
	"database/sql"
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

func handleAggregate(s *state, cmd command) error{
	if len(cmd.args) < 1{
		return fmt.Errorf("time between requests required")
	}

	//time_between_reqes := cmd.args[0]

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil{
		return err
	}
	fmt.Println("%v", feed)

	return nil
}

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

func scrapeFeeds(s *state) error{
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil{
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ID:  feed.ID, LastFetchedAt: sql.NullTime{time.Now(), true}})

	
	return nil
}