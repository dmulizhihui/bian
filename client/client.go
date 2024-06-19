package ba_client

import (
	"log"
	"net/http"
	"os"
)

type doFunc func(req *http.Request) (*http.Response, error)

// Client define API client
type Client struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
}

// Create client function for initialising new Binance client
func NewClient(apiKey string, secretKey string, baseURL ...string) *Client {
	url := "https://api.binance.com"

	if len(baseURL) > 0 {
		url = baseURL[0]
	}

	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    url,
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "ba_connect", log.LstdFlags),
	}
}
