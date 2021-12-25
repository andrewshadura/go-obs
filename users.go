// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"fmt"
	"net/http"
)

type WatchlistEntry struct {
	Name string `xml:"name,attr"`
}

type UserOwner struct {
	ID string `xml:"userid,attr"`
}

type User struct {
	ID        string           `xml:"login"`
	Email     string           `xml:"email"`
	Realname  string           `xml:"realname"`
	State     string           `xml:"state"`
	Owner     *UserOwner       `xml:"owner,omitempty"`
	Roles     []string         `xml:"globalrole,omitempty"`
	Watchlist []WatchlistEntry `xml:"watchlist>project"`
}

type person struct {
	User
}

type collection struct {
	Users []User `xml:"person"`
}

func (c *Client) GetUser(name string) (*User, error) {
	req, err := c.NewRequest(http.MethodGet, "/person/"+name, nil)
	if err != nil {
		return nil, err
	}

	var u person
	_, err = c.Do(req, &u)
	if err != nil {
		return nil, err
	}

	return &u.User, nil
}

type UserOptions struct {
	Prefix string `url:"prefix,omitempty"`
}

func (c *Client) GetUsers(prefix string) ([]string, error) {
	req, err := c.NewRequest(http.MethodGet, "/person", UserOptions{Prefix: prefix})
	if err != nil {
		return nil, err
	}

	var dir directory
	_, err = c.Do(req, &dir)
	if err != nil {
		return nil, err
	}

	var users []string
	for _, u := range dir.Entries {
		users = append(users, u.Name)
	}

	return users, nil
}

type SearchOptions struct {
	Match string `url:"match,omitempty"`
}

func (c *Client) GetUsersByEmail(email string) ([]User, error) {
	match := XPathAttrEquals("email", email).String()
	req, err := c.NewRequest(http.MethodGet, "/search/person", SearchOptions{Match: match})
	if err != nil {
		return nil, err
	}

	var results collection
	_, err = c.Do(req, &results)
	if err != nil {
		return nil, err
	}

	return results.Users, nil
}

func (c *Client) GetUserByEmail(email string) (*User, error) {
	users, err := c.GetUsersByEmail(email)
	if err != nil {
		return nil, err
	}

	if len(users) < 1 {
		return nil, fmt.Errorf("no users found with email %s", email)
	}

	if len(users) > 1 {
		return nil, fmt.Errorf("found %d users with email %s", len(users), email)
	}

	return &users[0], nil
}
