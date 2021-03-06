package jira

import (
	"fmt"
	"testing"
)

func TestBuildingSearchUrl(t *testing.T) {
	config := ClientConfig{
		Project:   "MyJiraProject",
		Subdomain: "companyname",
	}
	client := Client{Config: config}
	url := client.issueSearchURL("JQL")
	queryString := "jql=JQL&maxResults=500"
	expected := fmt.Sprintf("https://companyname.atlassian.net/rest/api/2/search?%s", queryString)
	if url != expected {
		t.Error("Expected valid JIRA URL: ", url)
	}
}

func TestAuthorizationHeaderConstruction(t *testing.T) {
	config := ClientConfig{
		Username: "Admin",
		Password: "Secret",
	}
	client := Client{Config: config}
	header := client.authorizationHeader()
	if header != "Basic QWRtaW46U2VjcmV0" {
		t.Error("Expected Valid Auth Header:", header)
	}
}

var apiResponse = []byte(`{"total":2,"issues":[{"id":"103220","fields": {"customfield_15000":null,"created":"2015-08-19T13:23:00.000-0600","resolutiondate":"2015-08-20T16:13:50.000-0600"}},{"id":"99686","fields":{"customfield_15000":null,"created":"2015-07-20T12:43:42.000-0600","resolutiondate":"2015-08-20T16:39:03.000-0600"}}]}`)

type stubAPI struct{}

func (api stubAPI) Fetch(url string, headers map[string]string) ([]byte, error) {
	return apiResponse, nil
}

func TestTranslatingJsonResponseToSearchResult(t *testing.T) {
	config := ClientConfig{
		Project:   "MyJiraProject",
		Subdomain: "companyname",
		Username:  "Admin",
		Password:  "Secret",
	}
	client := Client{Config: config}
	var stub stubAPI
	var result SearchResult
	client.QueryRecentlyClosedIssues(stub, "default", &result)
	if result.Total != 2 {
		t.Error("JSON not parsed correctly:", result)
	}
}
