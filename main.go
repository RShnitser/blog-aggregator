package main

import(
	"blog-aggregator/internal/config"
	"fmt"
)

func main(){
	cfg, err := config.Read()
	if err != nil{
		fmt.Printf("Error reading config: %v", err)
		return;
	}
	fmt.Println(cfg)
	
	err = cfg.SetUser("test")
	if err != nil{
		fmt.Printf("Error writing to config: %v", err)
		return;
	}
	fmt.Println(cfg)
}