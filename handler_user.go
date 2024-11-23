package main

import(
	"fmt"
	"context"
	"blog-aggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

func handlerRegister(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("username required")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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