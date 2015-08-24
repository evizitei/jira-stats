package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
)

// ClientConfig is just a data bag for stuff that comes out of the Configuration
// file.  The structure is entirely flat, and it's actual contents should never
//  be committed to a repo ever for any reason.
type ClientConfig struct {
	Username  string
	Password  string
	Subdomain string
	Project   string
}

// Client is the object that knows how to build intelligent API calls to JIRA,
// and how to deserialize responses into types defined in the "jira.go" file.
type Client struct {
	Config ClientConfig
}

var apiBasePath = "/rest/api/2"

func (c *Client) jqlBuilder(dateRange string) *Jql {
	return &Jql{
		Project:   c.Config.Project,
		DateRange: dateRange,
	}
}

func (c *Client) jiraHostURL() (*url.URL, error) {
	host := fmt.Sprintf("https://%s.atlassian.net", c.Config.Subdomain)
	jiraURL, err := url.Parse(host)
	if err != nil {
		println("URL Parsing Error:", err.Error)
		return nil, err
	}
	return jiraURL, nil
}

func (c *Client) issueSearchURL(jql string) string {
	jiraURL, _ := c.jiraHostURL()
	jiraURL.Path += fmt.Sprintf("%s/search", apiBasePath)
	parameters := url.Values{}
	parameters.Add("jql", jql)
	parameters.Add("maxResults", "500")
	jiraURL.RawQuery = parameters.Encode()
	return jiraURL.String()
}

func (c *Client) issueChangelogURL(issueKey string) string {
	jiraURL, _ := c.jiraHostURL()
	jiraURL.Path += fmt.Sprintf("%s/issue/%s", apiBasePath, issueKey)
	parameters := url.Values{}
	parameters.Add("expand", "changelog")
	jiraURL.RawQuery = parameters.Encode()
	return jiraURL.String()
}

func (c *Client) authorizationHeader() string {
	authString := fmt.Sprintf("%s:%s", c.Config.Username, c.Config.Password)
	encodedString := base64.StdEncoding.EncodeToString([]byte(authString))
	return fmt.Sprintf("Basic %s", encodedString)
}

func (c *Client) apiHeaders() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": c.authorizationHeader(),
	}
}

func (c *Client) queryIssues(api Api, result *SearchResult, jql string) error {
	searchURL := c.issueSearchURL(jql)
	responseBody, err := api.Fetch(searchURL, c.apiHeaders())
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

// QueryRecentlyClosedIssues will find all the JIRA tickets closed in the last 7 days
//  (or in the date range of your choice) and pack them into a search result object.
func (c *Client) QueryRecentlyClosedIssues(api Api, dateRange string, result *SearchResult) error {
	jql := c.jqlBuilder(dateRange).Closed()
	return c.queryIssues(api, result, jql)
}

// QueryRecentlyCreatedIssues will find all the JIRA tickets created in the last 7 days
//  (or in the date range of your choice) and pack them into a search result object.
func (c *Client) QueryRecentlyCreatedIssues(api Api, dateRange string, result *SearchResult) error {
	jql := c.jqlBuilder(dateRange).Created()
	return c.queryIssues(api, result, jql)
}

// QueryRecentlyDeployedIssues will find all the JIRA tickets deployed in the last 6 weeks
//  (or in the date range of your choice) and pack them into a search result object.  "Deployed"
//  in this case means having the "deployed-production" label applied, and is a convention.  This
//  should probably get expanded at some point to allow a more generic "label applied recently" query.
func (c *Client) QueryRecentlyDeployedIssues(api Api, dateRange string, result *SearchResult) error {
	jql := c.jqlBuilder(dateRange).Deployed()
	return c.queryIssues(api, result, jql)
}

// QueryChangelogsForResultSet is an expensive operation.  It takes a search result
//  set and iterates over every item, making an API call for each, in order to get
//  the detailed changelog history for each ticket.  It's sad, but some things that
//  you might want to query are only available in the ticket history.
func (c *Client) QueryChangelogsForResultSet(api Api, result *SearchResult) ([]*IssueHistory, error) {
	changelogs := make([]*IssueHistory, len(result.Issues))
	for i, issue := range result.Issues {
		issueURL := c.issueChangelogURL(issue.Key)
		responseBody, err := api.Fetch(issueURL, c.apiHeaders())
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
