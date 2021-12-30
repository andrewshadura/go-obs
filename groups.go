// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
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

type GroupMember struct {
	Username string `xml:"userid,attr"`
}

type Group struct {
	ID         string        `xml:"title"`
	Email      string        `xml:"email,omitempty"`
	Maintainer GroupMember   `xml:"maintainer"`
	Members    []GroupMember `xml:"person>person"`
}

func (g GroupMember) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if g.Username == "" {
		return nil
	} else {
		start.Attr = []xml.Attr{{
			Name: xml.Name{
				Space: "",
				Local: "userid",
			},
			Value: g.Username,
		}}
		return e.EncodeElement("", start)
	}
}

type group struct {
	Group
}

type directoryEntry struct {
	Name string `xml:"name,attr"`
}

type directory struct {
	Entries []directoryEntry `xml:"entry"`
}

func (c *Client) GetGroups() ([]string, error) {
	req, err := c.NewRequest(http.MethodGet, "/group", nil)
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

func (c *Client) GetGroup(name string) (*Group, error) {
	req, err := c.NewRequest(http.MethodGet, "/group/"+name, nil)
	if err != nil {
		return nil, err
	}

	var g group
	_, err = c.Do(req, &g)
	if err != nil {
		return nil, err
	}

	return &g.Group, nil
}

func (c *Client) NewGroup(name string) error {
	newGroup := group{Group{
		ID: name,
	}}
	req, err := c.NewRequest(http.MethodPut, "/group/"+name, newGroup)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteGroup(name string) error {
	req, err := c.NewRequest(http.MethodDelete, "/group/"+name, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateGroup(g *Group) error {
	name := g.ID
	gg := group{*g}
	req, err := c.NewRequest(http.MethodPut, "/group/"+name, gg)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AddGroupMember(group string, user string) error {
	req, err := c.NewRequest(http.MethodPost, "/group/"+group, UserOptions{Command: commandAddUser, User: user})
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveGroupMember(group string, user string) error {
	req, err := c.NewRequest(http.MethodPost, "/group/"+group, UserOptions{Command: commandRemoveUser, User: user})
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetGroupEmail(group string, email string) error {
	req, err := c.NewRequest(http.MethodPost, "/group/"+group, UserOptions{Command: commandSetEmail, Email: email})
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}
