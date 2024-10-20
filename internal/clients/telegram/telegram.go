package telegram

import (
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	httpClient *http.Client
	host       string
	path       string
}

func NewClient(host, token string) *Client {
	return &Client{
		host:       host,
		httpClient: &http.Client{},
		path:       path(token),
	}
}

func path(token string) string {
	return "bot" + token
}

func (c *Client) ReceiveClaim(offset, limit int) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

}

func (c *Client) SendAnswer() {

}
