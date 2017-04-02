package place

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	timeURL = "https://oauth.reddit.com/api/place/time.json"
	drawURL = "https://oauth.reddit.com/api/place/draw.json"
)

type Client struct {
	*http.Client
}

func (c Client) WaitTime() (time.Duration, error) {
	resp, err := c.Get(timeURL)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	var t Time
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return 0, err
	}

	return t.Duration(), nil
}

type Time struct {
	WaitSeconds float64 `json:"wait_seconds"`
}

func (t Time) Duration() time.Duration {
	d, _ := time.ParseDuration(fmt.Sprintf("%fs", t.WaitSeconds))
	return d
}

func (c Client) Draw(x, y int, color color.Color) (time.Duration, error) {
	cIndex := DefaultColorPalette.Index(color)

	v := url.Values{}
	v.Add("x", strconv.Itoa(x))
	v.Add("y", strconv.Itoa(y))
	v.Add("color", strconv.Itoa(cIndex))

	resp, err := c.PostForm(drawURL, v)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	var t Time
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return 0, err
	}

	return t.Duration(), nil
}
