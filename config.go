package main

import (
	"code.google.com/p/gcfg"
	"github.com/evizitei/jira-stats/jira"
)

type Config struct {
	Jira jira.JiraClientConfig
}

func LoadConfig() (Config, error) {
	var config Config
	err := gcfg.ReadFileInto(&config, "jirastats.gcfg")
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
