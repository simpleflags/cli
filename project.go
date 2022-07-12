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

type projectCommand struct {
	Account     string `short:"a" long:"acc" description:"Account identifier" env:"SF_ACCOUNT"`
	Name        string `short:"n" long:"name" description:"Project name"`
	Description string `short:"d" long:"description" description:"Describe your project"`
	Remove      bool   `long:"rm" description:"Remove project"`
	Args        struct {
		Identifier string `positional-arg-name:"identifier"`
	} `positional-args:"yes"`
}

func (c projectCommand) Execute(_ []string) error {
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

func (c projectCommand) create(ctx context.Context) error {
	if c.Account == "" {
		return errors.New("-a or --acc flag is required")
	}

	if c.Name == "" {
		return errors.New("-n or --name flag is required")
	}

	return api.CreateProject(ctx, &model.Project{
		Account:     c.Account,
		Identifier:  c.Args.Identifier,
		Name:        c.Name,
		Description: c.Description,
	})
}

func (c projectCommand) list(ctx context.Context) error {
	var accPtr *string
	if c.Account != "" {
		accPtr = &c.Account
	}
	projects, err := api.GetProjects(ctx, accPtr)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	columns := table.Row{"Identifier", "Name", "Description", "Account"}
	t.AppendHeader(columns)
	for _, val := range projects {
		item := table.Row{
			val.Identifier,
			val.Name,
			val.Description,
			val.Account,
		}

		t.AppendRow(item)
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}

func (c projectCommand) remove(ctx context.Context) error {
	if c.Account == "" {
		return errors.New("-a or --acc flag is required")
	}

	return api.DeleteProject(ctx, c.Account, c.Args.Identifier)
}

func init() {
	pc := projectCommand{}
	_, err := parser.AddCommand(
		"project",
		"Project commands",
		"Project based commands",
		&pc,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
