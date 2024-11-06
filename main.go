package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"log"
	"os"
	"database/sql"

	"github.com/hf-chow/gator/internal/command"
	"github.com/hf-chow/gator/internal/config"
	"github.com/hf-chow/gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	state := &command.State{Config: &cfg}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	dbQueries := database.New(db)
	state.DB = dbQueries

	cmds := &command.Commands{}
	cmds.Register("addfeed", command.HandlerAddFeed)
	cmds.Register("agg", command.HandlerAggregate)
	cmds.Register("feeds", command.HandlerFeed)
	cmds.Register("login", command.HandlerLogin)
	cmds.Register("register", command.HandlerRegister)
	cmds.Register("reset", command.HandlerReset)
	cmds.Register("users", command.HandlerUsers)

	args := os.Args
//	if (args[1] == "reset") || (args[1] == "users") || (args[1] == "agg") || (args[1] == "feeds") {
//		commandName := args[1]
//		commandArg := ""
//		cmd := command.Command{Name: commandName, Args: []string{commandArg}}
//		err := cmds.Run(state, cmd)
//		if err != nil {
//			fmt.Printf("Error %s\n", err)
//			os.Exit(1)
//		}
//		os.Exit(0)
//	}
	if len(args) < 2 {
		fmt.Println("Invalid input")
		os.Exit(1)
	}
	commandName := args[1]
	commandArg := args[2:]
	cmd := command.Command{Name: commandName, Args:commandArg}
	err = cmds.Run(state, cmd)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
