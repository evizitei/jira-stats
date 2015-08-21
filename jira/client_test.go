package jira

import (
  "testing"
  "fmt"
)

func TestJqlInterpolation(t *testing.T) {
  config := JiraClientConfig{
    Project: "MyJiraProject",
  }
  client := Client{ Config: config }
  jql := client.RecentlyClosedJql()
  if jql != "project=MyJiraProject AND status=Closed AND resolutiondate >= -7d" {
    t.Error("Expected correct jql interpolation: ", jql)
  }
}

func TestBuildingSearchUrl(t *testing.T) {
  config := JiraClientConfig{
    Project: "MyJiraProject",
    Subdomain: "companyname",
  }
  client := Client{ Config: config }
  url := client.IssueSearchUrl()
  queryString := "jql=project%3DMyJiraProject+AND+status%3DClosed+AND+resolutiondate+%3E%3D+-7d&maxResults=500"
  expected := fmt.Sprintf("https://companyname.atlassian.net/rest/api/2/search?%s", queryString)
  if url != expected {
    t.Error("Expected valid JIRA URL: ", url)
  }
}
