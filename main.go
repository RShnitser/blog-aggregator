package main

import(
	"blog-aggregator/internal/database"
	"database/sql"
	"blog-aggregator/internal/config"
	"fmt"
	"os"
	_ "github.com/lib/pq"
)

func main(){
	cfg, err := config.Read()
	if err != nil{
		fmt.Printf("Error reading config: %v", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil{
		fmt.Printf("Error creating database: %v", err)
		os.Exit(1)
	}
	
	s := state{database.New(db), &cfg}
	commands := commands{
		make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handleResetDatabase)
	commands.register("users", handleListUsers)
	commands.register("agg", handleAggregate)
	commands.register("addfeed", handleAddFeed)
	commands.register("feeds", handleListFeeds)

	args := os.Args[1:]
	if len(args) == 0{
		fmt.Println("no command provided")
		os.Exit(1)
	}

	command := command{
		args[0],
		args[1:],
	}

	err = commands.run(&s, command)
	if err != nil{
		fmt.Printf("could not execute command: %s\n", err)
		os.Exit(1)
	}
}