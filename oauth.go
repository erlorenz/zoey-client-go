package main

import (
	"errors"
	"net/http"

	"github.com/mrjones/oauth"
)

func (c *Client) newOauthClient(cfg Config) (*http.Client, error) {
	// Todo: set option for custom HTTP client

	// Create an oauth consumer with an option for custom http client
	consumer := oauth.NewConsumer(cfg.ConsumerKey, cfg.ConsumerSecret, oauth.ServiceProvider{})
	accessToken := &oauth.AccessToken{Token: cfg.TokenKey, Secret: cfg.TokenSecret}

	// Create an http client from the consumer and access token
	client, err := consumer.MakeHttpClient(accessToken)
	if err != nil {
		return nil, errors.New("failed to create new oauth client")
	}

	return client, nil
}
