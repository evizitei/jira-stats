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

var apiBasePath string = "/rest/api/2"

func (c *Client) buildJqlForProject(whereClause string) string {
	baseJql := "project=%s AND %s"
	project := c.Config.Project
	return fmt.Sprintf(baseJql, project, whereClause)
}

func (c *Client) RecentlyClosedJql() string {
	return c.buildJqlForProject("status=Closed AND resolutiondate >= -7d")
}

func (c *Client) RecentlyCreatedJql() string {
	return c.buildJqlForProject("createdDate >= -7d")
}

func (c *Client) RecentlyDeployedJql() string {
	return c.buildJqlForProject("labels in (deployed-production) AND resolutiondate >= -42d")
}

func (c *Client) jiraHostUrl() (*url.URL, error) {
	host := fmt.Sprintf("https://%s.atlassian.net", c.Config.Subdomain)
	jiraUrl, err := url.Parse(host)
	if err != nil {
		println("URL Parsing Error:", err.Error)
		return nil, err
	}
	return jiraUrl, nil
}

func (c *Client) IssueSearchUrl(jql string) string {
	jiraUrl, _ := c.jiraHostUrl()
	jiraUrl.Path += fmt.Sprintf("%s/search", apiBasePath)
	parameters := url.Values{}
	parameters.Add("jql", jql)
	parameters.Add("maxResults", "500")
	jiraUrl.RawQuery = parameters.Encode()
	return jiraUrl.String()
}

func (c *Client) IssueChangelogUrl(issueKey string) string {
	jiraUrl, _ := c.jiraHostUrl()
	jiraUrl.Path += fmt.Sprintf("%s/issue/%s", apiBasePath, issueKey)
	parameters := url.Values{}
	parameters.Add("expand", "changelog")
	jiraUrl.RawQuery = parameters.Encode()
	return jiraUrl.String()
}

func (c *Client) AuthorizationHeader() string {
	authString := fmt.Sprintf("%s:%s", c.Config.Username, c.Config.Password)
	encodedString := base64.StdEncoding.EncodeToString([]byte(authString))
	return fmt.Sprintf("Basic %s", encodedString)
}

func (c *Client) apiHeaders() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": c.AuthorizationHeader(),
	}
}

func (c *Client) queryIssues(api Api, result *SearchResult, jql string) error {
	searchUrl := c.IssueSearchUrl(jql)
	responseBody, err := api.Fetch(searchUrl, c.apiHeaders())
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

func (c *Client) QueryRecentlyClosedIssues(api Api, result *SearchResult) error {
	return c.queryIssues(api, result, c.RecentlyClosedJql())
}

func (c *Client) QueryRecentlyCreatedIssues(api Api, result *SearchResult) error {
	return c.queryIssues(api, result, c.RecentlyCreatedJql())
}

func (c *Client) QueryRecentlyDeployedIssues(api Api, result *SearchResult) error {
	return c.queryIssues(api, result, c.RecentlyDeployedJql())
}

func (c *Client) QueryChangelogsForResultSet(api Api, result *SearchResult) ([]*IssueHistory, error) {
	changelogs := make([]*IssueHistory, len(result.Issues))
	for i, issue := range result.Issues {
		issueUrl := c.IssueChangelogUrl(issue.Key)
		responseBody, err := api.Fetch(issueUrl, c.apiHeaders())
		if err != nil {
			println("Could not get result:", err.Error())
		}
		var history IssueHistory
		jsonErr := json.Unmarshal(responseBody, &history)
		if jsonErr != nil {
			println("Failed to parse JSON response:", jsonErr.Error())
		}
		changelogs[i] = &history
	}
	return changelogs, nil
}
