package main

import (
	"context"
	"github.com/antonmedv/expr"
	"github.com/kr/pretty"
	"github.com/simpleflags/services/pkg/api/client"
	"log"
	"os"
	"time"
)

type evaluateCommand struct {
	Account string            `short:"a" long:"acc" description:"Account identifier" env:"SF_ACCOUNT"`
	Project string            `short:"p" long:"project" description:"Project identifier" env:"SF_PROJECT"`
	Target  map[string]string `short:"t" long:"target" description:"Target data <property:value>"`
	Args    struct {
		Identifiers []string `positional-arg-name:"identifiers"`
	} `positional-args:"yes"`
}

func (c evaluateCommand) Execute(_ []string) error {
	clientAPI := client.New(os.Getenv("SF_API_KEY"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	target := make(map[string]any)
	for key, val := range c.Target {
		exprValue, err := expr.Eval(val, nil)
		target[key] = val
		if err == nil {
			target[key] = exprValue
		}
	}
	evaluate, err := clientAPI.Evaluate(ctx, nil, nil, nil, target, c.Args.Identifiers...)
	if err != nil {
		return err
	}
	_, err = pretty.Print(jsonFormatter("", "  ", evaluate))
	return err
}

func init() {
	ec := evaluateCommand{}
	cmd, err := parser.AddCommand(
		"eval",
		"Evaluate flag",
		"Evaluate flag with/out target",
		&ec,
	)

	cmd.Aliases = []string{"evaluate"}

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
