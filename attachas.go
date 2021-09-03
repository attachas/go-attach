package attachas

import (
	as "code.as/core/api"
)

const (
	apiURL    = "https://attach.as/api"
	devAPIURL = "https://dev.attach.as/api"
)

// Client is used to interact with the Attach.as API.
type Client struct {
	as.Client
}

// NewClient creates a new API client. All requests must be authenticated, so
// you should supply a user access token returned from the Write.as API
// library (github.com/writeas/go-writeas/v2)
//
//     wc := writeas.NewClient()
//     u, err := wc.LogIn("username", "password")
//     if err != nil {
//         // Handle error...
//     }
//     sc := attachas.NewClient(u.AccessToken)
func NewClient(token string) *Client {
	cfg := as.NewClientConfig(apiURL, "go-attachas v1")
	return NewClientWith(cfg, token)
}

// NewDevClient creates a new API client for development and testing. It will
// communicate with our development servers, and SHOULD NOT be used in
// production.
func NewDevClient(token string) *Client {
	return NewClientWith(as.NewClientConfig(devAPIURL, "go-attachas v1"), token)
}

// NewClientWith builds a new API client with the provided configuration.
func NewClientWith(cfg *as.ClientConfig, token string) *Client {
	cl := as.NewClient(cfg)
	cl.Token = token
	return &Client{*cl}
}
