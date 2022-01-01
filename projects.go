// Copyright (C) 2022, Andrej Shadura
// Copyright (C) 2022, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

// ProjectRef represents a project referred by its name.
// This is used e.g. to represent a project in a watchlist.
type ProjectRef struct {
	Name string `xml:"name,attr" json:"name"`
}
