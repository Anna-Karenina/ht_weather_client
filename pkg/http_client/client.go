package http_client

// Wrapper aroud go native http client

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	baseUrl string
	client  *http.Client
}

func NewClient(client *http.Client, base string) *Client {
	return &Client{
		client:  client,
		baseUrl: base,
	}
}

func (c *Client) RequestIntercepter(r *http.Request) (*http.Response, error) {
	r.Header.Add("Content-Type", "application/json")
	return c.client.Do(r)
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.RequestIntercepter(req)
}

func (c *Client) Post(url string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	return c.RequestIntercepter(req)
}

func (c *Client) Put(url string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	return c.RequestIntercepter(req)
}

func (c *Client) Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.RequestIntercepter(req)
}

func (c *Client) Path(url string) string {
	base := strings.TrimRight(c.baseUrl, "/")
	if url == "" {
		return base
	}
	return base + "/" + strings.TrimLeft(url, "/")
}

func (c *Client) Pathf(url string, v ...any) string {
	url = fmt.Sprintf(url, v...)
	return c.Path(url)
}
