package main

import(
	"strconv"
	"context"
	"fmt"
)

func handleBrowse(s *state, cmd command)error{
	limit := 2

	if len(cmd.args) == 1{
		newLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil{
			return err
		}
		limit = newLimit
	}

	posts, err := s.db.GetPosts(context.Background(), int32(limit))
	if err != nil{
		return err
	}

	for _, post := range posts{
		fmt.Println(post.Title)
	}

	return nil
}