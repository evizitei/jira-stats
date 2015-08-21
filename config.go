package main

import "code.google.com/p/gcfg"

type Config struct {
	Jira struct {
		Username string
		Password string
    Subdomain string
    Project string
	}
}

func LoadConfig() (Config, error) {
	var config Config
	err := gcfg.ReadFileInto(&config, "jirastats.gcfg")
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
