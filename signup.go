package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/simpleflags/cli/ui"
	"github.com/simpleflags/services/pkg/model"
	"log"
	"time"
)

type signupCommand struct {
}

func (c signupCommand) Execute(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	email := ui.EmailInput(ui.PromptContent{
		ErrMessage: "Please provide email",
		Label:      "Email",
	})

	password := ui.Password(ui.PromptContent{
		ErrMessage: "Password must have more than 6 characters",
	})

	confirmPassword := ui.Password(ui.PromptContent{
		ErrMessage: "Password must have more than 6 characters",
	})

	if password != confirmPassword {
		return errors.New("passwords mismatch")
	}

	err := api.Signup(ctx, &model.SignupBody{
		Email:          email,
		Password:       password,
		RepeatPassword: confirmPassword,
	})
	if err != nil {
		return err
	}
	fmt.Printf("User with email %s successfully registered\n", email)
	return nil
}

func init() {
	rc := signupCommand{}
	cmd, err := parser.AddCommand(
		"signup",
		"Register a new user to the SF system",
		"Type username and password to register to SF system",
		&rc,
	)

	cmd.Aliases = []string{"register"}

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
