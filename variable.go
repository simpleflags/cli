package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/jedib0t/go-pretty/table"
	"github.com/kr/pretty"
	"github.com/simpleflags/services/pkg/model"
	"log"
	"os"
	"time"
)

type variableCommand struct {
	Account     string            `short:"a" long:"acc" description:"Account identifier" env:"SF_ACCOUNT"`
	Project     *string           `short:"p" long:"project" description:"Project identifier" env:"SF_PROJECT"`
	Env         *string           `short:"e" long:"env" description:"Environment identifier used only to list variables"`
	Description string            `short:"d" long:"description" description:"Provide description for this variable"`
	Value       map[string]string `short:"v" long:"value" description:"Provide value for the variable in format <environment>:<value>"`
	Global      bool              `short:"g" long:"global" description:"Set this variable global"`
	Remove      bool              `long:"rm" description:"Remove variable"`
	Args        struct {
		Identifier string `positional-arg-name:"identifier"`
	} `positional-args:"yes"`
}

func (c *variableCommand) Execute(_ []string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// create or update
	if c.Global {
		c.Project = nil
	}

	if c.Args.Identifier == "" {
		return c.list(ctx)
	} else {
		if c.Remove {
			return c.remove(ctx)
		}
		if c.Env == nil && c.Description == "" && len(c.Value) == 0 {
			return c.print(ctx)
		}

		variables, err := api.GetVariables(ctx, c.Account, c.Project, c.Args.Identifier)
		if err != nil {
			return err
		}
		if len(variables) > 0 {
			return c.patch(ctx)
		}
		return c.create(ctx)
	}
}

func (c *variableCommand) list(ctx context.Context) error {
	variables, err := api.GetVariables(ctx, c.Account, c.Project)
	if err != nil {
		return err
	}
	columns := table.Row{"Account", "Project", "Identifier"}
	if c.Env != nil {
		columns = append(columns, "Value")
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(columns)
	for _, val := range variables {
		item := table.Row{val.Account, *val.Project, val.Identifier}
		if c.Env != nil {
			item = append(item, val.Value[*c.Env])
		}
		t.AppendRow(item)
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}

func (c *variableCommand) create(ctx context.Context) error {
	environments, err := api.GetEnvironments(ctx, &c.Account)
	if err != nil {
		return err
	}
	if len(environments) == 0 {
		return errors.New("no environments defined")
	}

	values := make(map[string]any)
	for key, val := range c.Value {
		value, err := expr.Eval(val, nil)
		values[key] = val
		if err == nil {
			values[key] = value
		}
	}

	body := model.Variable{
		Account:     c.Account,
		Project:     c.Project,
		Identifier:  c.Args.Identifier,
		Description: c.Description,
		Value:       values,
	}

	if err := api.CreateVariable(ctx, &body); err != nil {
		return err
	}
	project := ""
	if c.Project != nil {
		project = *c.Project
	}
	fmt.Printf("Variable %s successfully created in project %s", c.Args.Identifier,
		project)
	return nil
}

func (c *variableCommand) patch(ctx context.Context) error {
	for env, val := range c.Value {
		value, err := expr.Eval(val, nil)
		var newValue interface{}
		newValue = val
		if err == nil {
			newValue = value
		}
		body := model.PatchVariable{
			Value: newValue,
		}
		if err = api.PatchVariable(ctx, c.Account, c.Project, env, c.Args.Identifier, &body); err != nil {
			return err
		}
	}
	return nil
}

func (c *variableCommand) print(ctx context.Context) error {
	variables, err := api.GetVariables(ctx, c.Account, c.Project, c.Args.Identifier)
	if err != nil {
		return err
	}
	if len(variables) > 0 {
		_, err := pretty.Print(jsonFormatter("", "  ", variables[0]))
		return err
	}
	return nil
}

func (c *variableCommand) remove(ctx context.Context) error {
	err := api.DeleteVariable(ctx, c.Account, c.Project, c.Args.Identifier)
	if err != nil {
		return err
	}
	fmt.Printf("Variable %s successfully removed from project %s", c.Args.Identifier, c.Project)
	return nil
}

func init() {
	var err error
	vc := variableCommand{}
	cmd, err := parser.AddCommand(
		"var",
		"Variable operations",
		"Basic operations with variables",
		&vc,
	)

	cmd.Aliases = []string{"variable"}

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
