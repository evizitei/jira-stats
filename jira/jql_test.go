package jira

import (
	"strings"
	"testing"
)

func TestJqlInterpolation(t *testing.T) {
	jqlBuilder := Jql{Project: "MyJiraProject", DateRange: "default"}
	jql := jqlBuilder.Closed()
	if jql != "project=MyJiraProject AND status=Closed AND resolutiondate >= -7d" {
		t.Error("Expected correct jql interpolation: ", jql)
	}
}

func TestDefaultDateRange(t *testing.T) {
	jqlBuilder := Jql{Project: "Canvas", DateRange: "default"}
	jql := jqlBuilder.Created()
	if !strings.Contains(jql, "createdDate >= -7d") {
		t.Error("Expected default date range to be -7 days", jql)
	}
}

func TestDateRangeOverride(t *testing.T) {
	jqlBuilder := Jql{Project: "Canvas", DateRange: "2015-08-09:2015-08-22"}
	jql := jqlBuilder.Created()
	if !strings.Contains(jql, "createdDate >= 2015-08-09 AND createdDate <= 2015-08-22") {
		t.Error("Expected default date range to be overridden", jql)
	}
}

func TestDeployedRespectsOverride(t *testing.T) {
	jqlBuilder := Jql{Project: "Canvas", DateRange: "2015-08-09:2015-08-22"}
	jql := jqlBuilder.Deployed()
	expected := "project=Canvas AND labels in (deployed-production) AND resolutiondate >= 2015-08-09 AND resolutiondate <= 2015-08-22"
	if jql != expected {
		t.Error("Expected JQL for deployed tickets to respect date range", jql)
	}
}

func TestClosedRespectsOverride(t *testing.T) {
	jqlBuilder := Jql{Project: "Canvas", DateRange: "2015-08-09:2015-08-22"}
	jql := jqlBuilder.Closed()
	expected := "project=Canvas AND status=Closed AND resolutiondate >= 2015-08-09 AND resolutiondate <= 2015-08-22"
	if jql != expected {
		t.Error("Expected JQL for closed tickets to respect date range", jql)
	}
}
