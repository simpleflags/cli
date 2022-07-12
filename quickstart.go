package main

import "log"

type quickStartCommand struct {
}

func (c quickStartCommand) Execute(args []string) error {
	return nil
}

func init() {
	sc := quickStartCommand{}
	_, err := parser.AddCommand(
		"quick",
		"quick start template project",
		"create quickstart project",
		&sc,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
