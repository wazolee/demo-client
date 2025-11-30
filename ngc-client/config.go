// config.go
package main

import (
	"fmt"
	"strings"

	"ngc-client/crypto"
)

type AppConfig struct {
	Environment string
	UUID        string
	Hostname    string
	JwtToken    string            // environment alapú tokenek
	AuthBaseURL string
	UserBaseURL string
}

var configMap = map[string]AppConfig{
	"dev": {
		JwtToken:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.DEV_TOKEN_ITT",
		AuthBaseURL: "https://auth-dev.example.com",
		UserBaseURL: "https://user-dev.example.com",
	},
	"staging": {
		JwtToken:    "staging-jwt-token-123...",
		AuthBaseURL: "https://auth-staging.example.com",
		UserBaseURL: "https://user-staging.example.com",
	},
	"prod": {
		JwtToken:    "prod-jwt-token-very-secret...",
		AuthBaseURL: "https://auth.example.com",
		UserBaseURL: "https://user.example.com",
	},
}

func LoadConfig(encFile, keyHex string) (*AppConfig, error) {
	plain, err := crypto.DecryptFile(encFile, keyHex)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(strings.TrimSpace(plain), ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("érvénytelen formátum, várt: env:uuid:fqdn, kapott: %s", plain)
	}

	env := parts[0]
	cfg, exists := configMap[env]
	if !exists {
		return nil, fmt.Errorf("ismeretlen környezet: %s", env)
	}

	cfg.Environment = env
	cfg.UUID = parts[1]
	cfg.Hostname = parts[2]

	fmt.Printf("Betöltve config | Env: %s | Host: %s | UUID: %s\n",
		cfg.Environment, cfg.Hostname, cfg.UUID)

	return &cfg, nil
}