package jira

import (
	"testing"
	"time"

	"github.com/simplereach/timeutils"
)

var start time.Time
var end1 time.Time
var end2 time.Time

func buildResult() SearchResult {
	start, _ = timeutils.ParseDateString("09:51:20.939152pm 2015-09-08")
	end1, _ = timeutils.ParseDateString("09:51:20.939152pm 2015-10-08")
	end2, _ = timeutils.ParseDateString("09:51:20.939152pm 2015-12-08")

	var searchResult = SearchResult{
		Total: 2,
		Issues: []Issue{
			Issue{
				ID: "12345",
				Field: IssueFields{
					Created:  timeutils.Time{Time: start},
					Resolved: timeutils.Time{Time: end1},
				},
			},
			Issue{
				ID: "54321",
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
		ID: "12345",
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

func TestBugRatioDoesNotRound(t *testing.T) {
	issueTypes := []string{"Bug", "Bug", "Bug", "New Feature", "New Feature"}
	ratio := bugRatio(issueTypes)
	if ratio != 1.5 {
		t.Error("expected ratio not to round before dividing", ratio)
	}
}

func TestAvgAndMaxDiscardsZeroResults(t *testing.T) {
	times := []float64{2.0, 4.0, 6.0, 0.0}
	avg, _ := averageAndMax(times)
	if avg != 4 {
		t.Error("expected to discard zeros", avg)
	}
}
