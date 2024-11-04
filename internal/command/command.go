package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hf-chow/gator/internal/config"
	"github.com/hf-chow/gator/internal/database"
)

type Command struct {
	Name		string
	Args 		[]string
}

type Commands struct {
	Names		map[string]func(*State, Command) error
}

type State struct {
	DB			*database.Queries
	Config 		*config.Config
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("please provide the username")
	}

	if usernameExists(s, cmd.Args[0]) {
		err := s.Config.SetUser(cmd.Args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Username %s has been set\n", cmd.Args[0])
		return nil
	} else {
		os.Exit(1)
		return nil
	}
}

func usernameExists(s *State, username string) bool {
	user, err := s.DB.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("Unable to find user %s in DB\n", username)
		return false
	}
	fmt.Printf("User %s exists", user)
	return true
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("Please provide a username")
	}
	name := cmd.Args[0]
	args := database.CreateUserParams{
		ID: uuid.New(), CreatedAt: time.Now(),
		UpdatedAt: time.Now(), Name: name,
	}

	user, err := s.DB.CreateUser(context.Background(), args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return err
	}
	s.Config.SetUser(name)
	fmt.Printf("User %s has been created\n", name)
	fmt.Println(user)

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.Names == nil {
		c.Names = make(map[string]func(*State, Command) error)
	}
	c.Names[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	if f, ok := c.Names[cmd.Name]; ok {
		return f(s, cmd)
	}
	return errors.New("Command not found")
}
