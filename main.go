package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/caarlos0/env/v8"
	"github.com/mrjones/oauth"
)

type config struct {
	Port           int    `env:"PORT" envDefault:"3000"`
	ApiUrl         string `env:"ZOEY_API_URL,required"`
	ConsumerKey    string `env:"ZOEY_CONSUMER_KEY,required"`
	ConsumerSecret string `env:"ZOEY_CONSUMER_SECRET,required"`
	TokenKey       string `env:"ZOEY_TOKEN_KEY,required"`
	TokenSecret    string `env:"ZOEY_TOKEN_SECRET,required"`
}

func main() {
	// Parse the environment variables
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%+v\n", cfg)

	path := cfg.ApiUrl + "/api/rest/companyAccounts/account?id=100"
	params := map[string]string{
		"id": "100",
	}

	consumer := oauth.NewConsumer(cfg.ConsumerKey, cfg.ConsumerSecret, oauth.ServiceProvider{})
	accessToken := &oauth.AccessToken{Token: cfg.TokenKey, Secret: cfg.TokenSecret}

	consumer.AdditionalHeaders = map[string][]string{
		"Accept": {"application/json"},
	}

	resp, err := consumer.Get(path, params, accessToken)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

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
