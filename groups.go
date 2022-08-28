// Copyright (C) 2021, 2022 Andrej Shadura
// Copyright (C) 2021, 2022 Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"encoding/xml"
	"net/http"
)

const (
	commandAddUser    = "add_user"
	commandRemoveUser = "remove_user"
	commandSetEmail   = "set_email"
)

// Group represents a named group of users.
type Group struct {
	XMLName    xml.Name  `xml:"group"           json:"-"`
	ID         string    `xml:"title"           json:"name"`
	Email      string    `xml:"email,omitempty" json:"email,omitempty"`
	Maintainer UserRef   `xml:"maintainer"      json:"maintainer"`
	Members    []UserRef `xml:"person>person"   json:"members"`
}

type directoryEntry struct {
	Name string `xml:"name,attr"`
}

type directory struct {
	Entries []directoryEntry `xml:"entry"`
}

// ListGroups gets a list of names of all groups.
// Use GetGroup to retrieve the details of each group.
func (c *Client) ListGroups() ([]string, error) {
	req, err := c.NewRequest(http.MethodGet, "/group", nil, nil)
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

// GetGroup retrieves the details of the group (maintainer, members etc).
func (c *Client) GetGroup(name string) (*Group, error) {
	req, err := c.NewRequest(http.MethodGet, "/group/"+name, nil, nil)
	if err != nil {
		return nil, err
	}

	var g Group
	_, err = c.Do(req, &g)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// NewGroup creates a new empty group.
func (c *Client) NewGroup(name string) error {
	newGroup := Group{
		ID: name,
	}
	req, err := c.NewRequest(http.MethodPut, "/group/"+name, nil, newGroup)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteGroup deletes a group of users.
// On some OBS versions the group must be empty before it can be deleted.
func (c *Client) DeleteGroup(name string) error {
	req, err := c.NewRequest(http.MethodDelete, "/group/"+name, nil, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdateGroup updates an existing group.
// The user calling this must have necessary access rights to be able to
// perform this call.
func (c *Client) UpdateGroup(g *Group) error {
	name := g.ID
	req, err := c.NewRequest(http.MethodPut, "/group/"+name, nil, g)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// AddGroupMember adds a user to a group.
func (c *Client) AddGroupMember(group string, user string) error {
	req, err := c.NewRequest(http.MethodPost, "/group/"+group, UserOptions{Command: commandAddUser, User: user}, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// RemoveGroupMember removes a user from a group.
func (c *Client) RemoveGroupMember(group string, user string) error {
	req, err := c.NewRequest(http.MethodPost, "/group/"+group, UserOptions{Command: commandRemoveUser, User: user}, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// SetGroupEmail sets the email address of a group.
func (c *Client) SetGroupEmail(group string, email string) error {
	req, err := c.NewRequest(http.MethodPost, "/group/"+group, UserOptions{Command: commandSetEmail, Email: email}, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}
