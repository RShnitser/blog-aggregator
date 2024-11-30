package main

import(
	"strconv"
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

	return nil
}