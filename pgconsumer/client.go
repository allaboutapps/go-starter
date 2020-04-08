package pgconsumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"allaboutapps.at/aw/go-mranftl-sample/pgtestpool"
	"github.com/friendsofgo/errors"
)

type Client struct {
	baseURL *url.URL
	client  *http.Client
	config  ClientConfig
}

func NewClient(config ClientConfig) (*Client, error) {
	c := &Client{
		baseURL: nil,
		client:  nil,
		config:  config,
	}

	if len(c.config.BaseURL) == 0 {
		c.config.BaseURL = "http://pgserve:8080/api"
	}

	if len(c.config.APIVersion) == 0 {
		c.config.APIVersion = "v1"
	}

	u, err := url.Parse(c.config.BaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse base URL")
	}

	c.baseURL = u.ResolveReference(&url.URL{Path: path.Join(u.Path, c.config.APIVersion)})

	c.client = &http.Client{}

	return c, nil
}

func DefaultClientFromEnv() (*Client, error) {
	return NewClient(DefaultClientConfigFromEnv())
}

func (c *Client) ResetAllTracking(ctx context.Context) error {
	req, err := c.newRequest(ctx, "DELETE", "/admin/templates", nil)
	if err != nil {
		return err
	}

	var msg string
	resp, err := c.do(req, &msg)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.Errorf("failed to reset all tracking: %v", msg)
	}

	return nil
}

func (c *Client) InitializeTemplate(ctx context.Context, hash string) (*pgtestpool.TemplateDatabase, error) {
	payload := map[string]string{"hash": hash}

	req, err := c.newRequest(ctx, "POST", "/templates", payload)
	if err != nil {
		return nil, err
	}

	template := new(pgtestpool.TemplateDatabase)
	resp, err := c.do(req, template)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return template, nil
	case http.StatusLocked:
		return nil, pgtestpool.ErrTemplateAlreadyInitialized
	case http.StatusServiceUnavailable:
		return nil, pgtestpool.ErrManagerNotReady
	default:
		return nil, errors.Errorf("received unexpected HTTP status %d (%s)", resp.StatusCode, resp.Status)
	}
}

func (c *Client) FinalizeTemplate(ctx context.Context, hash string) error {
	req, err := c.newRequest(ctx, "PUT", fmt.Sprintf("/templates/%s", hash), nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return pgtestpool.ErrTemplateNotFound
	case http.StatusServiceUnavailable:
		return pgtestpool.ErrManagerNotReady
	default:
		return errors.Errorf("received unexpected HTTP status %d (%s)", resp.StatusCode, resp.Status)
	}
}

func (c *Client) GetTestDatabase(ctx context.Context, hash string) (*pgtestpool.TestDatabase, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/templates/%s/tests", hash), nil)
	if err != nil {
		return nil, err
	}

	db := new(pgtestpool.TestDatabase)
	resp, err := c.do(req, db)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return db, nil
	case http.StatusNotFound:
		return nil, pgtestpool.ErrTemplateNotFound
	case http.StatusServiceUnavailable:
		return nil, pgtestpool.ErrManagerNotReady
	default:
		return nil, errors.Errorf("received unexpected HTTP status %d (%s)", resp.StatusCode, resp.Status)
	}
}

func (c *Client) ReturnTestDatabase(ctx context.Context, hash string, id int) error {
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/templates/%s/tests/%d", hash, id), nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return pgtestpool.ErrTemplateNotFound
	case http.StatusServiceUnavailable:
		return pgtestpool.ErrManagerNotReady
	default:
		return errors.Errorf("received unexpected HTTP status %d (%s)", resp.StatusCode, resp.Status)
	}
}

func (c *Client) newRequest(ctx context.Context, method string, endpoint string, body interface{}) (*http.Request, error) {
	u := c.baseURL.ResolveReference(&url.URL{Path: path.Join(c.baseURL.Path, endpoint)})

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, errors.Wrap(err, "failed to encode request payload")
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}

	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusNoContent {
		return resp, nil
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, err
	}

	return resp, err
}
