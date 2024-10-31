package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"github.com/hf-chow/gator/internal/command"
	"github.com/hf-chow/gator/internal/config"
	"github.com/hf-chow/gator/internal/database"
	"os"
	"database/sql"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	state := &command.State{Config: &cfg}

	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)
	state.DB = dbQueries

	cmds := &command.Commands{}
	cmds.Register("login", command.HandlerLogin)

	args := os.Args
	if len(args) < 3 {
		fmt.Println("Invalid input")
		os.Exit(1)
	}
	commandName := args[1]
	commandArg := args[2]
	cmd := command.Command{Name: commandName, Args:[]string{commandArg}}
	err = cmds.Run(state, cmd)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
