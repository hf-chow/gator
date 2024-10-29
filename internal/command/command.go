package command

import (
	"errors"
	"fmt"
	"github.com/hf-chow/gator/internal/config"
)

type command struct {
	name		string
	args 		[]string
}

type commands struct {
	names		map[string]func(*state, command) error
}

type state struct {
	Config 		*config.Config
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("please provide the username")
	}
	err := s.Config.SetUser(cmd.args[1])
	if err != nil {
		return err
	}
	fmt.Printf("Username %s has been set\n", cmd.args[1])
 
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.names == nil {
		c.names = make(map[string]func(*state, command) error)
	}
	c.names[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.names[cmd.name]; ok {
		return f(s, cmd)
	}
	return errors.New("command not found")
}
