package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
)

type JiraClientConfig struct {
	Username  string
	Password  string
	Subdomain string
	Project   string
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

func (c *Client) AuthorizationHeader() string {
	authString := fmt.Sprintf("%s:%s", c.Config.Username, c.Config.Password)
	encodedString := base64.StdEncoding.EncodeToString([]byte(authString))
	return fmt.Sprintf("Basic %s", encodedString)
}

func (c *Client) QueryRecentlyClosedIssues(api Api, result *SearchResult) error {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": c.AuthorizationHeader(),
	}
	responseBody, err := api.Fetch(c.IssueSearchUrl(), headers)
	println("REsponse body is", responseBody)
	if err != nil {
		println("Could not get result:", err.Error())
		return err
	}
	jsonErr := json.Unmarshal(responseBody, result)
	if jsonErr != nil {
		println("Failed to parse JSON response:", jsonErr.Error())
		return jsonErr
	}
	return nil
}
