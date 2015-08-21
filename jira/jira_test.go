package jira

import (
  "testing"
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
