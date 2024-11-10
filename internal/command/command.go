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
	"github.com/hf-chow/gator/internal/parser"
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

func HandlerAddFeed(s * State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Please provide a name of the feed and the url")
	}
	if len(cmd.Args) < 2 {
		return errors.New("Please provide a url")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	createArgs := database.CreateFeedParams{
		ID: uuid.New(), 
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	}

	feed, err := s.DB.CreateFeed(context.Background(), createArgs)
	if err != nil {
		return err
	}

	feed_id := createArgs.ID

	followArgs := database.CreateFeedFollowParams{
		ID: uuid.New(), CreatedAt: time.Now(),
		UpdatedAt: time.Now(), UserID: user.ID,
		FeedID: feed_id,
	}
	follow, err := s.DB.CreateFeedFollow(context.Background(), followArgs)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", feed)
	fmt.Printf("%s\n", follow)
	return nil
}


func HandlerAggregate(s * State, cmd Command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		errors.New("invalid usage")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return errors.New("Invalid durations")
	}

	fmt.Printf("Collecting feeds every %s...", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <- ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *State) {
	feed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("Could not fetch next feed %s", err)
	}
	fmt.Println("Fetching next feed")
	scrapeFeed(s.DB, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Printf("Could not mark feed %s as fetched: %v", feed.Name, err)
		return 
	}

	data, err := parser.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Printf("Could not collect feed %s: %v", feed.Name, err)
		return 
	}
	for _, item := range data.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	fmt.Printf("Feed %s collected, %v posts found", feed.Name, len(data.Channel.Item))
}

func HandlerFeed(s *State, cmd Command) error {
	feeds, err := s.DB.GetFeed(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("%v", feeds)
	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Please provide an url")
	}
	url := cmd.Args[0]
	feed_id, err := s.DB.GetFeedIDByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	args := database.CreateFeedFollowParams{
		ID: uuid.New(), CreatedAt: time.Now(),
		UpdatedAt: time.Now(), UserID: user.ID,
		FeedID: feed_id,
	}
	follow, err := s.DB.CreateFeedFollow(context.Background(), args)

	fmt.Printf("%s", follow)

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	followings, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, following := range followings {
		fmt.Printf("%v\n", following)
	}

	return nil
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

func HandlerReset(s *State, cmd Command) error {
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		fmt.Printf("Table delete unsuccessful: %s", err)
		return err
	}
	fmt.Println("Table delete successful")
	return err
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.Config.CurrentUsername {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func HandlerUnfollow(s * State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Please provide an url")
	}
	url := cmd.Args[0]
	args := database.DeleteFeedFollowParams{
		Url: url, UserID: user.ID,
	}
	unfollow, err := s.DB.DeleteFeedFollow(context.Background(), args)
	if err != nil {
		return err
	}
	fmt.Printf("%s", unfollow)
	return nil
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

