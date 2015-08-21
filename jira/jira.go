package jira

import (
	"github.com/simplereach/timeutils"
)

type IssueFields struct {
	Created  timeutils.Time `json:"created"`
	Resolved timeutils.Time `json:"resolutiondate"`
}

type Issue struct {
	Id    string      `json:"id"`
	Key   string      `json:"key"`
	Field IssueFields `json:"fields"`
}

type SearchResult struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}
