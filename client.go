// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// Heavily inspired by go-gitlab by:
// Copyright (C) 2021, Sander van Harmelen
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	baseURL *url.URL

	// Authentication type used to make API calls.
	authType authType

	// Username and password used for basic authentication.
	username, password string

	// User agent used when communicating with the OBS API.
	UserAgent string

	// Don’t verify server’s TLS certificates
	InsecureSkipVerify bool
}

type authType int

const (
	basicAuth authType = iota
	userAgent          = "go-obs-api/0"
)

// NewAPI returns a new OBS API client. To use API methods which
// require authentication, provide a valid username and password.
func NewClient(username, password string, options ...ClientOptionFunc) (*Client, error) {
	client, err := newClient(options...)
	if err != nil {
		return nil, err
	}

	client.authType = basicAuth
	client.username = username
	client.password = password

	return client, nil
}

func newClient(options ...ClientOptionFunc) (*Client, error) {
	c := &Client{UserAgent: userAgent}

	c.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify},
		},
	}

	// Apply any given client options.
	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// BaseURL returns a copy of the baseURL.
func (c *Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
}

// setBaseURL sets the base URL for API requests to a custom endpoint.
func (c *Client) setBaseURL(urlStr string) error {
	urlStr = strings.TrimSuffix(urlStr, "/")
	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	c.baseURL = baseURL
	return nil
}

// NewRequest creates an API request. A relative URL path can be provided in
// path, in which case it is resolved relative to the base URL of the Client.
// If specified, the value pointed to by body is XML-encoded and included as
// the request body.
func (c *Client) NewRequest(method, path string, opt interface{}, body interface{}) (*http.Request, error) {
	u := *c.baseURL

	u.Path = c.baseURL.Path + path

	reqHeaders := make(http.Header)
	reqHeaders.Set("accept", "application/xml")

	if c.UserAgent != "" {
		reqHeaders.Set("user-agent", c.UserAgent)
	}

	var bodyReader io.Reader
	if body != nil {
		switch body := body.(type) {
		case string:
			bodyReader = bytes.NewReader([]byte(body))

		case interface{}:
			xml, err := xml.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(xml)
		}
	}

	if opt != nil {
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	// Set the request specific headers.
	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// XML-decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	// Set the correct authentication header.
	switch c.authType {
	case basicAuth:
		if c.username != "" {
			req.SetBasicAuth(c.username, c.password)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		// Even though there was an error, we still return the response
		// in case the caller wants to inspect it further.
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = xml.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}
