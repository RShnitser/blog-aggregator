package main

import(
	"blog-aggregator/internal/config"
	"fmt"
)

func main(){
	cfg,_ := config.Read()
	name := "test"
	fmt.Println(cfg)
	cfg.SetUser(&name)
	fmt.Println(cfg)
}