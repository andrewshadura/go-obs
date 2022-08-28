// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// Heavily inspired by go-gitlab by:
// Copyright (C) 2021, Sander van Harmelen
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// An ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	Body     []byte         `xml:"-"`
	Response *http.Response `xml:"-"`
	Message  string         `xml:"summary"`
	Code     string         `xml:"code,attr"`
	XMLName  xml.Name       `xml:"status"`
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent, http.StatusNotModified:
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Body = data

		var raw ErrorResponse
		if err := xml.Unmarshal(data, &raw); err != nil {
			errorResponse.Message = fmt.Sprintf("failed to parse unknown error format: '%s'", data)
		} else {
			errorResponse.Message = raw.Message
			errorResponse.Code = raw.Code
		}
	}

	return errorResponse
}
