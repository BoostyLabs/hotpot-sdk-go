package examples

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/BoostyLabs/hotpot-sdk-go/client"
)

// Config represents the configurations for examples.
type Config struct {
	ApiKey  string `env:"HOTPOT_API_KEY,required,unset"`
	BaseUrl string `env:"HOTPOT_BASE_URL,unset" envDefault:"https://api.hotpot.tech,unset"`

	WalletAddresses []string  `env:"WALLET_ADDRESSES,unset" envSeparator:","`
	RetailID        string    `env:"RETAILER_ID,unset"`
	Token           string    `env:"TOKEN,unset"`
	TokenQuery      string    `env:"TOKEN_QUERY,unset"`
	NetworkID       int64     `env:"NETWORK_ID,unset"`
	ActiveSwaps     bool      `env:"ACTIVE_SWAPS,unset"`
	Limit           int64     `env:"LIMIT,unset" envDefault:"10"`
	Offset          int64     `env:"OFFSET,unset" envDefault:"0"`
	IntentID        uuid.UUID `env:"INTENT_ID,unset"`
}

// LoadConfig loads the configurations from environment variables.
func LoadConfig() *Config {
	envFile := os.Getenv("ENV_FILE")
	if envFile != "" {
		if err := godotenv.Overload(envFile); err != nil {
			log.Fatalf("failed to load config file: %v", err)
		}
	}

	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatalf("failed to parse env: %v", err)
	}

	return &config
}

// InitClient initializes the client from the config values.
func (c *Config) InitClient() *client.Client {
	return client.NewDefault(c.BaseUrl, c.ApiKey)
}
