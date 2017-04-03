package place

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const loginBaseURL = "https://www.reddit.com/api/login/"

func Login(username, password string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	transport := &uaTransport{}
	client := &http.Client{
		Transport: transport,
		Jar:       jar,
	}

	res, err := client.PostForm(loginBaseURL+username,
		url.Values{
			"op":       {"login"},
			"user":     {username},
			"passwd":   {password},
			"api_type": {"json"},
		})

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	defer res.Body.Close()

	var reply redditJSON
	if err := json.NewDecoder(res.Body).Decode(&reply); err != nil {
		return nil, err
	}

	if len(reply.JSON.Errors) != 0 {
		return nil, errors.New(fmt.Sprintf("error: %s", reply.JSON.Errors))
	}

	var dataJSON struct {
		ModHash string `json:"modhash"`
	}

	if err := json.Unmarshal(reply.JSON.Data, &dataJSON); err != nil {
		return nil, err
	}

	transport.modHash = dataJSON.ModHash
	log.Println(dataJSON.ModHash)

	log.Printf("[DEBUG] %s\n", reply)

	return &Client{client}, nil
}

type redditJSON struct {
	JSON struct {
		Errors []json.RawMessage `json:"errors"`
		Data   json.RawMessage   `json:"data"`
	} `json:"json"`
}

func main() {
}

type uaTransport struct {
	modHash string
}

func (t *uaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "go:vallon.me/place:v0.0.1 (by /u/neilvallon)")

	if t.modHash != "" {
		req.Header.Set("x-modhash", t.modHash)
	}

	resp, err := http.DefaultTransport.RoundTrip(req)

	if err == nil {
		log.Println("[INFO] X-Ratelimit-Used", resp.Header["X-Ratelimit-Used"])
		log.Println("[INFO] X-Ratelimit-Remaining", resp.Header["X-Ratelimit-Remaining"])
		log.Println("[INFO] X-Ratelimit-Reset", resp.Header["X-Ratelimit-Reset"])
	}

	return resp, err
}
