package main

import (
	"github.com/simpleflags/cli/config"
	"github.com/simpleflags/services/pkg/api/admin"
	"io/ioutil"
	"log"
	"path"
)

var (
	api *admin.API
)

func init() {
	sfDir, err := config.GetSimpleFlagsDir()
	if err != nil {
		log.Fatalf("Error getting data directory %v", err)
	}
	token, _ := ioutil.ReadFile(path.Join(sfDir, "auth.data"))
	api = admin.New(string(token))
}
