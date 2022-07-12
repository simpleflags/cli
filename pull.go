package main

import (
	"context"
	sfsdk "github.com/simpleflags/golang-server-sdk"
	"github.com/simpleflags/golang-server-sdk/client"
	"github.com/simpleflags/golang-server-sdk/connector/simple"
	"github.com/simpleflags/golang-server-sdk/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type pullCommand struct {
	Interval int      `short:"i" long:"interval" description:"Pull interval (60s)" default:"60"`
	Flags    []string `short:"f" long:"flags" description:"Flags"`
}

func (c pullCommand) Execute(_ []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fileStorage, err := repository.NewFileStorage("./")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	conn := simple.NewHttpConnector(os.Getenv("SF_API_KEY"),
		simple.WithBaseURL("https://64a55c46.fanoutcdn.com/api"))
	err = sfsdk.InitWithConnector(conn, client.WithStorage(&fileStorage))
	if err != nil {
		log.Printf("could not connect to SF servers %v", err)
	}

	defer func() {
		if err := sfsdk.Close(); err != nil {
			log.Printf("error while closing client err: %v", err)
		}
	}()
	sfsdk.WaitForInitialization()

	<-ctx.Done()
	return nil
}

func init() {
	sc := pullCommand{}
	_, err := parser.AddCommand(
		"pull",
		"Pull data from server and listen for changes",
		"Pull data from server and listen for changes",
		&sc,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
