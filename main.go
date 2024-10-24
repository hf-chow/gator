package main

import (
	"fmt"
	"github.com/hf-chow/gator/internal/config"
	"encoding/json"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%s", err)
	}

	username := "hon"
	err = cfg.SetUser(username)
	if err != nil {
		fmt.Printf("%s", err)
	}

	newCfg, err := config.Read()
	if err != nil {
		fmt.Printf("%s", err)
	}
	b, _ := json.Marshal(&newCfg)
	fmt.Printf("%s", string(b))
}
