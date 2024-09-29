package main

import(
	"blog-aggregator/internal/config"
	"fmt"
	"os"
)

func main(){
	cfg, err := config.Read()
	if err != nil{
		fmt.Printf("Error reading config: %v", err)
		os.Exit(1)
	}
	
	s := state{&cfg}
	commands := commands{
		make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)

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