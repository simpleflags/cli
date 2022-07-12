package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/simpleflags/services/pkg/model"
	"log"
	"os"
	"time"
)

type accountCommand struct {
	Remove bool `long:"rm" description:"Remove account"`
	Args   struct {
		Name string `positional-arg-name:"name"`
	} `positional-args:"yes"`
}

func (c accountCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if c.Args.Name != "" {
		return c.createAccount(ctx)
	}

	if c.Remove {
		return c.remove(ctx)
	}

	return c.list(ctx)
}

func (c accountCommand) list(ctx context.Context) error {
	accounts, err := api.GetAccounts(ctx)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	columns := table.Row{"Identifier", "Name", "Owner"}
	t.AppendHeader(columns)
	for _, val := range accounts {
		item := table.Row{
			val.Identifier,
			val.Name,
			val.Owner,
		}

		t.AppendRow(item)
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}

func (c accountCommand) remove(ctx context.Context) error {
	return api.DeleteAccount(ctx, c.Args.Name)
}

func (c accountCommand) createAccount(ctx context.Context) error {
	if c.Args.Name == "" {
		return errors.New("you need to provide name of account")
	}
	body := model.CreateAccountBody{
		Name: c.Args.Name,
	}

	response, err := api.CreateAccount(ctx, &body)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created account with identifier %s", response.Identifier)
	return nil
}

func init() {
	ac := accountCommand{}
	cmd, err := parser.AddCommand(
		"acc",
		"Account commands",
		"Account based commands",
		&ac,
	)

	cmd.Aliases = []string{"account"}

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
