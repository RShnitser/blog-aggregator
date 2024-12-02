package main

import(
	"blog-aggregator/internal/database"
	"strconv"
	"context"
	"fmt"
)

func handleBrowse(s *state, cmd command, user database.User)error{
	limit := 2

	if len(cmd.args) == 1{
		newLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil{
			return err
		}
		limit = newLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{ UserID: user.ID, Limit: int32(limit)})
	
	if err != nil{
		return err
	}

	for _, post := range posts{
		fmt.Println(post.Title)
	}

	return nil
}