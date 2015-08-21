package jira

import (
  "fmt"
  "net/url"
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

func (c *Client) IssueSearchUrl() string {
  host := fmt.Sprintf("https://%s.atlassian.net", c.Config.Subdomain)
  jiraUrl, err := url.Parse(host)
  if err != nil {
		println("URL Parsing Error:", err.Error)
		return ""
	}
  jiraUrl.Path += "/rest/api/2/search"
  parameters := url.Values{}
  parameters.Add("jql", c.RecentlyClosedJql())
  parameters.Add("maxResults", "500")
  jiraUrl.RawQuery = parameters.Encode()
  return jiraUrl.String()
}
