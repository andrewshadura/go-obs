// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
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

type WatchlistEntry struct {
	Name string `xml:"name,attr" json:"name"`
}

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

type User struct {
	XMLName   xml.Name         `xml:"person"               json:"-"`
	ID        string           `xml:"login"                json:"username"`
	Email     string           `xml:"email"                json:"email"`
	Realname  string           `xml:"realname"             json:"realname"`
	State     string           `xml:"state"                json:"state"`
	Owner     *UserRef         `xml:"owner,omitempty"      json:"owner,omitempty"`
	Roles     []string         `xml:"globalrole,omitempty" json:"globalrole,omitempty"`
	Watchlist []WatchlistEntry `xml:"watchlist>project"    json:"watchlist,omitempty"`
}

type collection struct {
	Users []User `xml:"person"`
}

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

func (c *Client) LockUser(name string) error {
	req, err := c.NewRequest(http.MethodPost, "/person/"+name, UserOptions{Command: commandLockUser, User: name})
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteUser(name string) error {
	req, err := c.NewRequest(http.MethodPost, "/person/"+name, UserOptions{Command: commandDeleteUser, User: name})
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}
