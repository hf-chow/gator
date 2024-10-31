package command

import (
	"errors"
	"fmt"

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
	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Username %s has been set\n", cmd.Args[0])
 
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
