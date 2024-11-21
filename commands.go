package main

import(
	"blog-aggregator/internal/database"
	"blog-aggregator/internal/config"
	"fmt"
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

