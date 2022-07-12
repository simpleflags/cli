package main

import (
	"context"
	"errors"
	"github.com/jedib0t/go-pretty/table"
	"github.com/simpleflags/services/pkg/model"
	"log"
	"os"
	"time"
)

type envCommand struct {
	Account     string `short:"a" long:"acc" description:"Account identifier" env:"SF_ACCOUNT"`
	Name        string `short:"n" long:"name" description:"Project name"`
	Description string `short:"d" long:"description" description:"Describe your project"`
	Production  bool   `long:"prod" description:"Production environment"`
	Remove      bool   `long:"rm" description:"Remove environment"`
	Args        struct {
		Identifier string `positional-arg-name:"identifier"`
	} `positional-args:"yes"`
}

func (c envCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if c.Args.Identifier == "" {
		return c.list(ctx)
	} else {
		if c.Remove {
			return c.remove(ctx)
		}

		return c.create(ctx)
	}
}

func (c envCommand) create(ctx context.Context) error {
	if c.Account == "" {
		return errors.New("-a or --acc flag is required")
	}

	if c.Name == "" {
		return errors.New("-n or --name flag is required")
	}

	return api.CreateEnvironment(ctx, &model.Environment{
		Account:     c.Account,
		Identifier:  c.Args.Identifier,
		Name:        c.Name,
		Description: c.Description,
		Production:  c.Production,
	})
}

func (c envCommand) list(ctx context.Context) error {
	var accPtr *string
	if c.Account != "" {
		accPtr = &c.Account
	}

	envs, err := api.GetEnvironments(ctx, accPtr)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	columns := table.Row{"Identifier", "Name", "Description", "Production", "Account"}
	t.AppendHeader(columns)
	for _, val := range envs {
		item := table.Row{
			val.Identifier,
			val.Name,
			val.Description,
			val.Production,
			val.Account,
		}

		t.AppendRow(item)
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}

func (c envCommand) remove(ctx context.Context) error {
	if c.Account == "" {
		return errors.New("-a or --acc flag is required")
	}

	return api.DeleteEnvironment(ctx, c.Account, c.Args.Identifier)
}

func init() {
	ec := envCommand{}
	_, err := parser.AddCommand(
		"env",
		"Environment commands",
		"Environment based commands",
		&ec,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
