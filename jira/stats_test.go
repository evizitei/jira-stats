package jira

import (
	"github.com/simplereach/timeutils"
	"testing"
	"time"
)

var start time.Time
var end1 time.Time
var end2 time.Time

func buildResult() SearchResult {
	start, _ = timeutils.ParseDateString("09:51:20.939152pm 2015-09-08")
	end1, _ = timeutils.ParseDateString("09:51:20.939152pm 2015-10-08")
	end2, _ = timeutils.ParseDateString("09:51:20.939152pm 2015-12-08")

	var searchResult SearchResult = SearchResult{
		Total: 2,
		Issues: []Issue{
			Issue{
				Id: "12345",
				Field: IssueFields{
					Created:  timeutils.Time{Time: start},
					Resolved: timeutils.Time{Time: end1},
				},
			},
			Issue{
				Id: "54321",
				Field: IssueFields{
					Created:  timeutils.Time{Time: start},
					Resolved: timeutils.Time{Time: end2},
				},
			},
		},
	}

	return searchResult
}

func TestCycleTimeAverageCalculation(t *testing.T) {
	avg, _ := CalculateCycleTime(buildResult())
	if int(avg) != 48 {
		t.Error("expected accurate average calculation", int(avg))
	}
}

func TestCycleTimeMaxCalculation(t *testing.T) {
	_, max := CalculateCycleTime(buildResult())
	if int(max) != 72 {
		t.Error("expected accurate max calculation", int(max))
	}
}

func TestCycleTimeCalculation(t *testing.T) {
	issue := Issue{
		Id: "12345",
		Field: IssueFields{
			Created:  timeutils.Time{Time: start},
			Resolved: timeutils.Time{Time: end1},
		},
	}
	cycleTime := cycleTimeForIssue(issue)
	if int(cycleTime) != 24 {
		t.Error("expected good duration check", int(cycleTime))
	}
}
