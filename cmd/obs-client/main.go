// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/andrewshadura/go-obs"
	"github.com/urfave/cli/v2"
	"github.com/zalando/go-keyring"
)

var client *obs.Client

func formatOutput(c *cli.Context, data interface{}) {
	if c.Bool("json") {
		jsonOutput, _ := json.MarshalIndent(data, "", "\t")
		fmt.Printf("%s\n", jsonOutput)
	} else {
		switch v := reflect.ValueOf(data); v.Kind() {
		case reflect.Array, reflect.Slice:
			switch v.Index(0).Kind() {
			case reflect.String:
				fmt.Println(strings.Join(v.Interface().([]string), "\n"))
			default:
				for i := 0; i < v.Len(); i++ {
					fmt.Printf("%#v\n", v.Index(i))
				}
			}
		case reflect.Interface, reflect.Ptr:
			fmt.Printf("%#v\n", v.Interface())
		}
	}
}

func userListCmd(c *cli.Context) error {
	prefix := ""
	if c.NArg() > 0 {
		prefix = c.Args().First()
	}

	users, err := client.ListUsers(prefix)
	if err != nil {
		return fmt.Errorf("failed to list users: %s", err)
	}

	formatOutput(c, users)

	return nil
}

func userGetCmd(c *cli.Context) error {
	user, err := client.GetUser(c.Args().First())
	if err != nil {
		return fmt.Errorf("failed to retrieve user: %s", err)
	}

	formatOutput(c, user)

	return nil
}

func userLookupCmd(c *cli.Context) error {
	user, err := client.GetUsersByEmail(c.Args().First())
	if err != nil {
		return fmt.Errorf("failed to look up user: %s", err)
	}

	formatOutput(c, user)

	return nil
}

func groupListCmd(c *cli.Context) error {
	groups, err := client.ListGroups()
	if err != nil {
		return fmt.Errorf("failed to retrieve groups: %s", err)
	}

	formatOutput(c, groups)

	return nil
}

func userGetGroupsCmd(c *cli.Context) error {
	groups, err := client.GetUserGroups(c.Args().First())
	if err != nil {
		return fmt.Errorf("failed to retrieve user's groups: %s", err)
	}

	formatOutput(c, groups)

	return nil
}

func groupGetCmd(c *cli.Context) error {
	group, err := client.GetGroup(c.Args().First())
	if err != nil {
		return fmt.Errorf("failed to retrieve group: %s", err)
	}

	formatOutput(c, group)

	return nil
}

func groupNewCmd(c *cli.Context) error {
	err := client.NewGroup(c.Args().First())
	if err != nil {
		return fmt.Errorf("failed to create group: %s", err)
	}

	group, err := client.GetGroup(c.Args().First())
	if err != nil {
		return fmt.Errorf("failed to retrieve group after its creation: %s", err)
	}

	formatOutput(c, group)

	return nil
}

func groupDeleteCmd(c *cli.Context) error {
	group, err := client.GetGroup(c.Args().First())
	if err != nil {
		return fmt.Errorf("failed to retrieve group: %s", err)
	}

	group.Members = nil

	err = client.UpdateGroup(group)
	if err != nil {
		return fmt.Errorf("failed to remove users from group: %s", err)
	}

	err = client.DeleteGroup(c.Args().First())
	if err != nil {
		return fmt.Errorf("failed to delete group: %s", err)
	}

	return nil
}

func groupAddCmd(c *cli.Context) error {
	err := client.AddGroupMember(c.Args().Get(1), c.Args().Get(0))
	if err != nil {
		return fmt.Errorf("failed to add user %s to group %s: %s", c.Args().Get(0), c.Args().Get(1), err)
	}

	return nil
}

func groupRemoveCmd(c *cli.Context) error {
	err := client.RemoveGroupMember(c.Args().Get(1), c.Args().Get(0))
	if err != nil {
		return fmt.Errorf("failed to remove user %s from group %s: %s", c.Args().Get(0), c.Args().Get(1), err)
	}

	return nil
}

type urlFlag struct {
	Url url.URL
}

func (u *urlFlag) Set(value string) error {
	parsed, err := url.Parse(value)
	if err == nil {
		u.Url = *parsed
	}

	return err
}

func (u *urlFlag) String() string {
	return u.Url.String()
}

func parseUrlFlag(value string) *urlFlag {
	u := urlFlag{}
	if u.Set(value) != nil {
		return nil
	} else {
		return &u
	}
}

func main() {
	app := &cli.App{
		Usage:                "OBS API command-line client",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.GenericFlag{
				Name:  "api-url",
				Value: parseUrlFlag("https://build.opensuse.org/"),
				Usage: "OBS API `URL` (including auth info)",
			},
			&cli.BoolFlag{
				Name:  "use-keyring",
				Value: true,
				Usage: "Use keyring for passwords",
			},
			&cli.BoolFlag{
				Name:  "json",
				Value: false,
				Usage: "Output results in JSON",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "user",
				Usage: "Manipulate users",
				Subcommands: []*cli.Command{
					{
						Name:      "list",
						Usage:     "List all users with username starting with PREFIX",
						Action:    userListCmd,
						ArgsUsage: "[PREFIX]",
					},
					{
						Name:      "get",
						Usage:     "Get a user by username",
						Action:    userGetCmd,
						ArgsUsage: "USERNAME",
					},
					{
						Name:      "lookup",
						Usage:     "Lookup users by email",
						Action:    userLookupCmd,
						ArgsUsage: "EMAIL",
					},
					{
						Name:      "groups",
						Usage:     "List groups of a user",
						Action:    userGetGroupsCmd,
						ArgsUsage: "USERNAME",
					},
				},
			},
			{
				Name:  "group",
				Usage: "Manipulate groups",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List all groups",
						Action: groupListCmd,
					},
					{
						Name:   "new",
						Usage:  "Create a new group",
						Action: groupNewCmd,
					},
					{
						Name:   "get",
						Usage:  "Get a group by its name",
						Action: groupGetCmd,
					},
					{
						Name:   "delete",
						Usage:  "Delete a group",
						Action: groupDeleteCmd,
					},
					{
						Name:   "add",
						Usage:  "Add a user to a group",
						Action: groupAddCmd,
					},
					{
						Name:   "remove",
						Usage:  "Remove a user from a group",
						Action: groupRemoveCmd,
					},
				},
			},
		},
		Before: func(c *cli.Context) error {
			if u, ok := c.Generic("api-url").(*urlFlag); ok {
				apiUrl := u.Url
				user := apiUrl.User.Username()
				var pass string
				if c.Bool("use-keyring") {
					pass, _ = keyring.Get(apiUrl.Host, user)
				}

				explicitPass, havePass := apiUrl.User.Password()
				if havePass {
					pass = explicitPass
				}
				apiUrl.User = nil

				var err error
				client, err = obs.NewClient(user, pass, obs.WithBaseURL(apiUrl.String()))
				if err != nil {
					log.Fatalf("failed to create client: %s", err)
				}
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
