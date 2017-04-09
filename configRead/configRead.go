// Package configRead reads configuration information from either a JSON file or from
// static environment variables as a fallback.
package configRead

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Config holds the different configuration options.
type Config struct {
	DiscordKey          string   `json:"discord-key"`
	DiscordSecret       string   `json:"discord-secret"`
	AuthURL             string   `json:"discord-auth-url"`
	HTTPPort            string   `json:"http-port"`
	Hostname            string   `json:"hostname,omitempty"`
	AuthRedirectHandler string   `json:"auth-redirect-handler"`
	BotServerName       string   `json:"bot-server-name"`
	BotServerPath       string   `json:"bot-server-path"`
	Scopes              []string `json:"scopes"`
}

// ReadConfig takes a path to a JSON file.
// If it fails to read the file, it falls back to environment variables.
// Returns an error if it can't parse the JSON file or if it can't read environment variables.
func ReadConfig(path string) (Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil { // not using json config. Try to get it from env vars
		config := Config{
			DiscordKey:          os.Getenv("DISCORD_KEY"),
			DiscordSecret:       os.Getenv("DISCORD_SECRET"),
			AuthURL:             os.Getenv("AUTH_URL"),
			HTTPPort:            os.Getenv("PORT"),
			Hostname:            os.Getenv("HOSTNAME"),
			AuthRedirectHandler: os.Getenv("AUTH_REDIRECT"),
			BotServerName:       os.Getenv("BOT_SERVER"),
			BotServerPath:       os.Getenv("BOT_SERVER_PATH"),
			Scopes:              []string{"identify"},
		}
		if !strings.Contains(config.HTTPPort, ":") {
			config.HTTPPort = ":" + config.HTTPPort
		}
		fmt.Println(config)
		return config, nil
	}
	config := Config{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}
	return config, err
}
