package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	ApiUrl         string `env:"ZOEY_API_URL,required"`
	ConsumerKey    string `env:"ZOEY_CONSUMER_KEY,required"`
	ConsumerSecret string `env:"ZOEY_CONSUMER_SECRET,required"`
	TokenKey       string `env:"ZOEY_TOKEN_KEY,required"`
	TokenSecret    string `env:"ZOEY_TOKEN_SECRET,required"`
}

type Client struct {
	consumerKey    string
	consumerSecret string
	tokenKey       string
	tokenSecret    string
	apiUrl         string
	Client         *http.Client
}

func NewClient(cfg Config) *Client {
	c := &Client{
		consumerKey:    cfg.ConsumerKey,
		consumerSecret: cfg.ConsumerSecret,
		tokenKey:       cfg.TokenKey,
		tokenSecret:    cfg.TokenSecret,
		apiUrl:         cfg.ApiUrl,
	}

	oc, err := c.newOauthClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create new client: %s", err)
	}

	c.Client = oc
	return c
}

// NewRequest mirrors http.NewRequest to return *http.Request.
// It adds the Accept and Content-Type headers and the ApiUrl from the config
func (c *Client) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.apiUrl+url, body)
	if err != nil {
		return nil, err
	}

	// Add headers to not get blocked by Cloudflare
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func main() {
	// Parse the environment variables
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%v\n", err)
	}

	// Instantiate the new client with credientials and API url
	zc := NewClient(cfg)

	req, err := zc.NewRequest("GET", "/api/rest/companyAccounts/account?id=822", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := zc.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {

		errorString := fmt.Sprintf("Request came back with failing status code: %s\n", resp.Status)
		log.Fatalf(errorString)
	}

	var body struct {
		Name        string
		External_id string
	}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("JSON: %+v\n", body)

}
