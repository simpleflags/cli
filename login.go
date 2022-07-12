package main

import (
	"context"
	"fmt"
	"github.com/simpleflags/cli/config"
	"github.com/simpleflags/cli/ui"
	"github.com/simpleflags/services/pkg/model"
	"io/ioutil"
	"log"
	"path"
	"time"
)

type loginCommand struct {
	Args struct {
		Email string `positional-arg-name:"email"`
	} `positional-args:"yes"`
}

func (c loginCommand) Execute(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	email := c.Args.Email
	if email == "" {
		email = ui.TextInput(ui.PromptContent{
			ErrMessage: "Please provide email",
			Label:      "Email",
		})
	}
	password := ui.Password(ui.PromptContent{
		ErrMessage: "Password must have more than 6 characters",
	})
	fmt.Println("Authenticating...")
	response, err := api.Authenticate(ctx, &model.LoginRequestBody{
		Email:    email,
		Password: password,
	})
	if err != nil {
		fmt.Printf("error while authenticating, err: %v", err)
		return err
	}

	sfDir, err := config.GetSimpleFlagsDir()
	if err != nil {
		fmt.Printf("error getting simpleflags folder, err: %v", err)
	}

	err = ioutil.WriteFile(path.Join(sfDir, "auth.data"), []byte(response.Token), 0644)
	if err != nil {
		fmt.Printf("error saving token, err: %v", err)
		return err
	}
	fmt.Printf("Welcome %s, how are you doing today?\n", email)
	return nil
}

func init() {
	lc := loginCommand{}
	_, err := parser.AddCommand(
		"login",
		"Login to the SF system",
		"Type username and password to login to SF system",
		&lc,
	)
	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
