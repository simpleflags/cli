package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/simpleflags/cli/config"
	"github.com/simpleflags/services/pkg/model"
	"log"
	"path"
	"time"
)

type apiKeyCommand struct {
	Account  string `short:"a" long:"acc" description:"Account identifier" env:"SF_ACCOUNT"`
	Project  string `short:"p" long:"project" description:"Project identifier" required:"true" env:"SF_PROJECT"`
	Env      string `short:"e" long:"env" description:"Environment identifier"`
	Name     string `short:"n" long:"name" description:"Key name" required:"true"`
	Remove   bool   `long:"rm" description:"Remove flag"`
	SetAsEnv bool   `long:"set-env" description:"Set api key as environment variable"`

	Permissions map[string]bool `long:"perm" description:"Set permission example create_account:true"`

	Args struct {
		Identifier string `positional-arg-name:"identifier"`
	} `positional-args:"yes"`
}

func (c apiKeyCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if c.Remove {
		return c.remove(ctx)
	}
	return c.create(ctx)
}

func (c apiKeyCommand) create(ctx context.Context) error {
	var permissions model.Permissions
	for key, val := range c.Permissions {
		c.accountPermissions(key, val, &permissions)
		c.projectPermissions(key, val, &permissions)
		c.envPermissions(key, val, &permissions)
		c.keyPermissions(key, val, &permissions)
		c.variablePermissions(key, val, &permissions)
		c.flagPermissions(key, val, &permissions)
	}

	response, err := api.CreateAPIKey(ctx, &model.APIKey{
		Account:     c.Account,
		Project:     c.Project,
		Environment: c.Env,
		Identifier:  c.Args.Identifier,
		Name:        c.Name,
		Permissions: permissions,
	})
	if err != nil {
		return err
	}

	fmt.Printf("API Key %s created\n", response.Key)
	if c.SetAsEnv {
		return c.setAsEnv(response.Key)
	}

	return nil
}

func (c apiKeyCommand) accountPermissions(key string, val bool, permissions *model.Permissions) {
	switch key {
	case "account":
		permissions.Account.Create = val
		permissions.Account.Update = val
		permissions.Account.List = val
		permissions.Account.View = val
		permissions.Account.Delete = val
	case "create_account", "account_create":
		permissions.Account.Create = val
	case "update_account", "account_update":
		permissions.Account.Update = val
	case "list_account", "list_accounts", "account_list":
		permissions.Account.List = val
	case "view_account", "account_view":
		permissions.Account.View = val
	case "delete_account", "account_delete":
		permissions.Account.Delete = val
	}
}

func (c apiKeyCommand) projectPermissions(key string, val bool, permissions *model.Permissions) {
	switch key {
	case "project":
		permissions.Project.Create = val
		permissions.Project.Update = val
		permissions.Project.List = val
		permissions.Project.View = val
		permissions.Project.Delete = val
	case "create_project", "project_create":
		permissions.Project.Create = val
	case "update_project", "project_update":
		permissions.Project.Update = val
	case "list_project", "list_projects", "project_list":
		permissions.Project.List = val
	case "view_project", "project_view":
		permissions.Project.View = val
	case "delete_project", "project_delete":
		permissions.Project.Delete = val
	}
}

func (c apiKeyCommand) envPermissions(key string, val bool, permissions *model.Permissions) {
	switch key {
	case "env", "environment":
		permissions.Environment.Create = val
		permissions.Environment.Update = val
		permissions.Environment.List = val
		permissions.Environment.View = val
		permissions.Environment.Delete = val
	case "create_env", "env_create":
		permissions.Environment.Create = val
	case "update_env", "env_update":
		permissions.Environment.Update = val
	case "list_env", "list_envs", "env_list":
		permissions.Environment.List = val
	case "view_env", "env_view":
		permissions.Environment.View = val
	case "delete_env", "env_delete":
		permissions.Environment.Delete = val
	}
}

func (c apiKeyCommand) keyPermissions(key string, val bool, permissions *model.Permissions) {
	switch key {
	case "key":
		permissions.Key.Create = val
		permissions.Key.Update = val
		permissions.Key.List = val
		permissions.Key.View = val
		permissions.Key.Delete = val
	case "create_key", "key_create":
		permissions.Key.Create = val
	case "update_key", "key_update":
		permissions.Key.Update = val
	case "list_key", "list_keys", "key_list":
		permissions.Key.List = val
	case "view_key", "key_view":
		permissions.Key.View = val
	case "delete_key", "key_delete":
		permissions.Key.Delete = val
	}
}

func (c apiKeyCommand) variablePermissions(key string, val bool, permissions *model.Permissions) {
	switch key {
	case "var", "variable":
		permissions.Variable.Create = val
		permissions.Variable.Update = val
		permissions.Variable.List = val
		permissions.Variable.View = val
		permissions.Variable.Delete = val
	case "create_var", "var_create", "create_variable", "variable_create":
		permissions.Variable.Create = val
	case "update_var", "var_update", "update_variable", "variable_update":
		permissions.Variable.Update = val
	case "list_var", "list_vars", "var_list", "list_variable", "list_variables", "variable_list":
		permissions.Variable.List = val
	case "view_var", "var_view", "view_variable", "variable_view":
		permissions.Variable.View = val
	case "delete_var", "var_delete", "delete_variable", "variable_delete":
		permissions.Variable.Delete = val
	}
}

func (c apiKeyCommand) flagPermissions(key string, val bool, permissions *model.Permissions) {
	switch key {
	case "flag":
		permissions.Flag.Create = val
		permissions.Flag.Update = val
		permissions.Flag.List = val
		permissions.Flag.View = val
		permissions.Flag.Delete = val
	case "create_flag", "flag_create":
		permissions.Flag.Create = val
	case "update_flag", "flag_update":
		permissions.Flag.Update = val
	case "list_flag", "list_flags", "flag_list":
		permissions.Flag.List = val
	case "view_flag", "flag_view":
		permissions.Flag.View = val
	case "delete_flag", "flag_delete":
		permissions.Flag.Delete = val
	}
}

func (c apiKeyCommand) remove(ctx context.Context) error {
	err := api.DeleteAPIKey(ctx, c.Account, c.Project, c.Env, c.Args.Identifier)
	if err != nil {
		return err
	}
	fmt.Printf("API key %s successfully removed from project %s and environment %s",
		c.Args.Identifier, c.Project, c.Env)
	return nil
}

func (c apiKeyCommand) setAsEnv(apiKey string) error {
	sfDir, err := config.GetSimpleFlagsDir()
	if err != nil {
		log.Fatalf("Error getting data directory %v", err)
	}

	envs, err := godotenv.Read(path.Join(sfDir, ".env"))
	if err != nil {
		log.Println(err)
	}

	envs["SF_API_KEY"] = apiKey

	err = godotenv.Write(envs, path.Join(sfDir, ".env"))
	if err == nil {
		fmt.Printf("API key stored as env variable SF_API_KEY")
	}
	return err
}

func init() {
	akc := apiKeyCommand{}
	_, err := parser.AddCommand(
		"key",
		"API Key commands",
		"adding and removing API Keys",
		&akc,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
