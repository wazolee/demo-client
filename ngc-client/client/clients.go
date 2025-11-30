package client

import (
	"net/http"
	"time"
)

type Config struct {
	Environment string `mapstructure:"environment"`
	Services    struct {
		Auth struct {
			BaseURL   string `mapstructure:"base_url"`
			JwtToken  string `mapstructure:"jwt_token"`
		} `mapstructure:"auth"`
		User struct {
			BaseURL string `mapstructure:"base_url"`
		} `mapstructure:"user"`
	} `mapstructure:"services"`
}

type Clients struct {
	AuthClient *http.Client
	UserClient *http.Client
	AuthToken  string
	AuthURL    string
	UserURL    string
}

func NewClients(cfg Config) (*Clients, error) {
	// Auth client – JWT tokennel
	authTransport := http.DefaultTransport.(*http.Transport).Clone()
	authClient := &http.Client{
		Timeout:   15 * time.Second,
		Transport: authTransport,
	}

	// User client – ugyanaz, de a token a headerben lesz hozzáadva híváskor
	userClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	return &Clients{
		AuthClient: authClient,
		UserClient: userClient,
		AuthToken:  cfg.Services.Auth.JwtToken,
		AuthURL:    cfg.Services.Auth.BaseURL,
		UserURL:    cfg.Services.User.BaseURL,
	}, nil
}

func (c *Clients) DoAuthRequest(method, path string) (*http.Response, error) {
	url := c.AuthURL + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	return c.AuthClient.Do(req)
}

func (c *Clients) DoUserRequest(method, path string) (*http.Response, error) {
	url := c.UserURL + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	return c.UserClient.Do(req)
}