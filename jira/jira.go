package jira

import (
	"github.com/simplereach/timeutils"
)

// IssueType represents the options for what a JIRA ticket is
// (bug, story, feature, task, etc)
type IssueType struct {
	Name string `json:"name"`
}

// IssueFields is a struct that kind of feels clunky because it's an attribute
//  map but one level down from it's parent.  Really, it represents the "fields"
//  attribute for Issues, which is there so that custom fields don't cause an
//  issue for JIRA's API JSON structure.  There are many more fields other than
//  those defined here, but these are the only ones we've needed so far
type IssueFields struct {
	Created   timeutils.Time `json:"created"`
	Resolved  timeutils.Time `json:"resolutiondate"`
	IssueType IssueType      `json:"issuetype"`
}

// Issue represents a ticket in JIRA.  The "Key" is the token by which we usually
//  refer to a ticket ([PROJECT_ABBREVIATION]-[SEQUENTIAL_ID]).  All the fields
//  you actually care about will be found in the "Field" object.
type Issue struct {
	ID    string      `json:"id"`
	Key   string      `json:"key"`
	Field IssueFields `json:"fields"`
}

type SearchResult struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}

type IssueChangeItem struct {
	Field     string `json:"field"`
	FieldType string `json:"fieldtype"`
	From      string `json:"fromString"`
	To        string `json:"toString"`
}

type IssueChange struct {
	Id      string            `json:"id"`
	Created timeutils.Time    `json:"created"`
	Items   []IssueChangeItem `json:"items"`
}

type Changelog struct {
	Total     int           `json:"total"`
	Histories []IssueChange `json:"histories"`
}

type IssueHistory struct {
	Id        string    `json:"id"`
	Key       string    `json:"key"`
	Changelog Changelog `json:"changelog"`
}
