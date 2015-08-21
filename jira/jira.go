package jira

import (
  "fmt"
  "github.com/simplereach/timeutils"
)

type JiraClientConfig struct {
  Username string
  Password string
  Subdomain string
  Project string
}

type Client struct {
  Config JiraClientConfig
}

func (c *Client) RecentlyClosedJql() string {
  baseJql := "project=%s AND status=Closed AND resolutiondate >= -7d"
  project := c.Config.Project
  return fmt.Sprintf(baseJql, project)
}

type IssueFields struct {
  Created timeutils.Time `json:"created"`
  Resolved timeutils.Time `json:"resolutiondate"`
}

type Issue struct {
  Id string `json:"id"`
  Field IssueFields `json:"fields"`
}

type SearchResult struct {
	Total     int    `json:"total"`
	Issues    []Issue `json:"issues"`
}
