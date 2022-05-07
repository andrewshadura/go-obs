// Copyright (C) 2021, 2022 Andrej Shadura
// Copyright (C) 2021, 2022 Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

const (
	commandChangePassword = "change_password" //nolint
	commandLockUser       = "lock"
	commandDeleteUser     = "delete"
)

type UserOptions struct {
	Command string `url:"cmd,omitempty"`
	User    string `url:"userid,omitempty"`
	Email   string `url:"email,omitempty"`
	Prefix  string `url:"prefix,omitempty"`
}

// UserRef represents a user referred by their username.
// This is used e.g. to represent members of a group.
type UserRef struct {
	ID string `xml:"userid,attr" json:"username"`
}

func (u UserRef) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if u.ID == "" {
		return nil
	} else {
		start.Attr = []xml.Attr{{
			Name: xml.Name{
				Space: "",
				Local: "userid",
			},
			Value: u.ID,
		}}
		return e.EncodeElement("", start)
	}
}

func (u UserRef) MarshalJSON() ([]byte, error) {
	if u.ID == "" {
		return json.Marshal(nil)
	} else {
		return json.Marshal(u.ID)
	}
}

// User represents a user (a person in OBS terminology)
type User struct {
	XMLName   xml.Name     `xml:"person"               json:"-"`
	ID        string       `xml:"login"                json:"username"`
	Email     string       `xml:"email"                json:"email"`
	Realname  string       `xml:"realname"             json:"realname"`
	State     string       `xml:"state"                json:"state"`
	Owner     *UserRef     `xml:"owner,omitempty"      json:"owner,omitempty"`
	Roles     []string     `xml:"globalrole,omitempty" json:"globalrole,omitempty"`
	Watchlist []ProjectRef `xml:"watchlist>project"    json:"watchlist,omitempty"`
}

type collection struct {
	Users []User `xml:"person"`
}

// GetUser retrieves the details of the user (email, real name etc).
func (c *Client) GetUser(name string) (*User, error) {
	req, err := c.NewRequest(http.MethodGet, "/person/"+name, nil)
	if err != nil {
		return nil, err
	}

	var u User
	_, err = c.Do(req, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// ListUsers gets a list of names of all users.
// Use GetUser to retrieve the details of each user.
func (c *Client) ListUsers(prefix string) ([]string, error) {
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

// GetUsersByEmail returns the details of the users matching given email address
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

// GetUserByEmail returns the details of the only user matching given email address
// If more than one user with given email address exist, an error is returned.
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

// LockUser locks the user and their projects
func (c *Client) LockUser(name string) error {
	req, err := c.NewRequest(http.MethodPost, "/person/"+name, UserOptions{Command: commandLockUser})
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser marks the user as deleted and deletes their projects
func (c *Client) DeleteUser(name string) error {
	req, err := c.NewRequest(http.MethodPost, "/person/"+name, UserOptions{Command: commandDeleteUser}, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetUserGroups retrieves a list of groups the user is a member of
func (c *Client) GetUserGroups(name string) ([]string, error) {
	req, err := c.NewRequest(http.MethodGet, "/person/"+name+"/group", nil)
	if err != nil {
		return nil, err
	}

	var dir directory
	_, err = c.Do(req, &dir)
	if err != nil {
		return nil, err
	}

	var groups []string
	for _, g := range dir.Entries {
		groups = append(groups, g.Name)
	}

	return groups, nil
}
