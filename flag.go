package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/jedib0t/go-pretty/table"
	"github.com/kr/pretty"
	"github.com/simpleflags/evaluation"
	"github.com/simpleflags/services/pkg/model"
	"log"
	"os"
	"time"
)

type flagCommand struct {
	Account     string              `short:"a" long:"acc" description:"Account identifier" env:"SF_ACCOUNT"`
	Project     string              `short:"p" long:"project" description:"Project identifier" required:"true" env:"SF_PROJECT"`
	Env         string              `short:"e" long:"env" description:"Environment identifier (use only when modifying rules)"`
	Name        string              `short:"n" long:"name" description:"Flag name"`
	Description string              `short:"d" long:"description" description:"Provide description for this flag"`
	Permanent   *bool               `long:"permanent" description:"Permanent flag"`
	Deprecated  *bool               `long:"deprecated" description:"Deprecated flag"`
	On          *bool               `long:"on" description:"Flag activation"`
	Off         *bool               `long:"off" description:"Flag deactivate"`
	OffValue    string              `long:"off-value" description:"Provide off value"`
	Rules       []map[string]string `short:"r" long:"rule" description:"Provide rule expression for the value"`
	Tags        []string            `short:"t" long:"tag" description:"Tags"`
	Remove      bool                `long:"rm" description:"Remove flag"`
	Args        struct {
		Identifier string `positional-arg-name:"identifier"`
	} `positional-args:"yes"`
}

func (c flagCommand) Execute(_ []string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if c.Args.Identifier == "" {
		return c.list(ctx)
	} else {
		if c.Remove {
			return c.removeFlag(ctx)
		}

		if c.Name == "" && c.Description == "" && c.On == nil && c.Off == nil &&
			c.OffValue == "" && len(c.Rules) == 0 && len(c.Tags) == 0 {
			return c.print(ctx)
		}
		return c.createOrUpdate(ctx)
	}
}

func (c flagCommand) print(ctx context.Context) error {
	flag, err := api.GetFlag(ctx, c.Account, c.Project, c.Args.Identifier)
	if err != nil {
		return err
	}
	_, err = pretty.Print(jsonFormatter("", "  ", flag))
	return err
}

func (c flagCommand) list(ctx context.Context) error {
	flags, err := api.GetFlags(ctx, c.Account, c.Project)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	columns := table.Row{"Project", "Name", "Identifier", "Permanent", "Deprecated", "Version"}
	t.AppendHeader(columns)
	for _, val := range flags {
		item := table.Row{
			val.Project,
			val.Name,
			val.Identifier,
			val.Permanent,
			val.Deprecated,
			val.Version,
		}

		t.AppendRow(item)
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}

func (c flagCommand) createOrUpdate(ctx context.Context) error {
	flag, err := api.GetFlag(ctx, c.Account, c.Project, c.Args.Identifier)
	if err == nil && flag.Identifier != "" {
		return c.patch(ctx)
	}

	tags, err := api.GetTags(ctx, c.Account, c.Project, c.Args.Identifier)
	if err == nil && len(tags) > 0 {
		return c.patch(ctx)
	}

	return c.create(ctx)
}

func (c flagCommand) patch(ctx context.Context) error {
	var instructions model.Instructions
	if c.Name != "" {
		instructions.Name = c.Name
	}

	if c.Description != "" {
		instructions.Description = c.Description
	}

	instructions.Permanent = c.Permanent
	instructions.Deprecated = c.Deprecated

	if len(c.Tags) > 0 {
		instructions.AddTags = c.Tags
	}

	// environment based patch

	if c.On != nil || c.Off != nil || c.OffValue != "" || len(c.Rules) > 0 {
		if c.Env == "" {
			return errors.New("environment -e or --env flag is required")
		}

		if c.OffValue != "" {
			instructions.SetOffValue.Value = c.OffValue
			eval, err := expr.Eval(c.OffValue, nil)
			if err == nil {
				instructions.SetOffValue.Value = eval
			} else {
				instructions.SetOffValue.Value = c.OffValue
			}
			instructions.SetOffValue.Environment = c.Env
		}

		if c.On != nil && c.Off != nil {
			return errors.New("cannot set on and off flags in same command")
		}

		if c.On != nil {
			instructions.SetOn.Value = true
			instructions.SetOn.Environment = c.Env
		}

		if c.Off != nil {
			instructions.SetOn.Value = false
			instructions.SetOn.Environment = c.Env
		}

		if len(c.Rules) > 0 {
			for _, rule := range c.Rules {
				for e, val := range rule {
					eval, err := expr.Eval(val, nil)
					var value interface{} = val
					if err == nil {
						value = eval
					}
					instructions.Rules = append(instructions.Rules, model.RuleInstruction{
						Environment: c.Env,
						Value:       value,
						Expression:  e,
					})
				}
			}
		}
	}

	return api.PatchFlag(ctx, c.Account, c.Project, c.Args.Identifier, &instructions)
}

func (c flagCommand) create(ctx context.Context) error {
	if c.Name == "" {
		return errors.New("please provide a flag name with -n or --name option")
	}
	environments, err := api.GetEnvironments(ctx, &c.Account)
	if err != nil {
		return err
	}
	if len(environments) == 0 {
		return errors.New("no environments defined")
	}

	offValue, err := expr.Eval(c.OffValue, nil)
	if err != nil {
		offValue = c.OffValue
	}

	rules := make([]evaluation.Rule, len(c.Rules))
	for i, rule := range c.Rules {
		for e, val := range rule {
			eval, err := expr.Eval(val, nil)
			var value interface{} = val
			if err == nil {
				value = eval
			}
			rules[i] = evaluation.Rule{
				Expression: e,
				Value:      value,
			}
		}
	}

	envs := make(map[string]model.Configuration)
	for _, env := range environments {
		envs[env.Identifier] = model.Configuration{
			OffValue: offValue,
			Rules:    rules,
		}
	}
	permanent := false
	if c.Permanent != nil {
		permanent = *c.Permanent
	}

	body := model.CreateFlagBody{
		Account:      c.Account,
		Project:      c.Project,
		Identifier:   c.Args.Identifier,
		Name:         c.Name,
		Description:  &c.Description,
		Permanent:    permanent,
		Environments: envs,
		Tags:         c.Tags,
	}
	defer func() {
		if err == nil {
			fmt.Printf("Flag %s successfully created in project %s", c.Args.Identifier,
				c.Project)
		}
	}()
	return api.CreateFlag(ctx, &body)
}

func (c flagCommand) removeFlag(ctx context.Context) error {
	err := api.DeleteFlag(ctx, c.Account, c.Project, c.Args.Identifier)
	if err != nil {
		return err
	}
	fmt.Printf("Flag %s successfully deleted in project %s", c.Args.Identifier, c.Project)
	return nil
}

func init() {
	var err error
	cfc := flagCommand{}
	_, err = parser.AddCommand(
		"flag",
		"Flag commands",
		"Create, update and list flags",
		&cfc,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
