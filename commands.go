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

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
	})

	return err
}