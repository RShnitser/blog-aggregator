package main

import(
	"context"
	"fmt"
)

func handleResetDatabase(s *state, cmd command) error{
	err := s.db.DeleteUsers(context.Background())
	if err != nil{
		return err
	}
	fmt.Println("All users deleted")
	return nil
}