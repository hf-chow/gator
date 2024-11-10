package command

import (
	"context"

	"github.com/hf-chow/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUsername)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}

}
