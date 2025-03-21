package xkcd

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"log/slog"
	"net/http"
	"time"

	"yadro.com/course/update/core"
)

type XKCDResponse struct {
	Num        int    `json:"num"`
	Transcript string `json:"transcript"`
	Title      string `json:"title"`
}

type Client struct {
	log    *slog.Logger
	client http.Client
	url    string
}

func NewClient(url string, timeout time.Duration, log *slog.Logger) (*Client, error) {
	if url == "" {
		return nil, fmt.Errorf("empty base url specified")
	}
	if timeout < 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}
	return &Client{
		client: http.Client{Timeout: timeout},
		log:    log,
		url:    url,
	}, nil
}

func (c Client) Get(_ context.Context, id int) (core.XKCDInfo, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", c.url, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error("Fail to create req", "err", err)
		return core.XKCDInfo{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("Error while getting XKCD", "error", err)
		return core.XKCDInfo{}, err
	}
	defer core.SafeClose(resp.Body)

	var result XKCDResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.log.Error("Failed to decode XKCD response", "error", err)
		return core.XKCDInfo{}, err
	}

	res := core.XKCDInfo{
		ID:          result.Num,
		URL:         fmt.Sprintf("%s/%d/info.0.json", c.url, result.Num),
		Title:       result.Title,
		Description: result.Transcript,
	}
	return res, nil
}

func (c Client) LastID(_ context.Context) (int, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/info.0.json", c.url))
	if err != nil {
		c.log.Error("Error while getting last id from xkcd", "error", err)
		return -1, err
	}
	defer core.SafeClose(resp.Body)

	var result XKCDResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.log.Error("Failed to decode XKCD response", "error", err)
		return -1, err
	}
	return result.Num, nil
}
