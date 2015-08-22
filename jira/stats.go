package jira

import (
	"time"
	"strings"
	"errors"
)

func cycleTimeForIssue(issue Issue) float64 {
	resolved := issue.Field.Resolved
	duration := resolved.Sub(issue.Field.Created.Time)
	return duration.Hours()
}

func findInProgressTime(changes []IssueChange) (time.Time, error) {
  for _, issueChange := range changes {
		for _, changeItem := range issueChange.Items {
			if changeItem.Field == "status" && changeItem.To == "In Progress" {
				return issueChange.Created.Time, nil
			}
		}
	}
	return time.Now(), errors.New("No In Progress Time Found...")
}

func findDeployTime(changes []IssueChange) (time.Time, error) {
	for _, issueChange := range changes {
		for _, changeItem := range issueChange.Items {
			if changeItem.Field == "labels"{
				if strings.Contains(changeItem.To, "deployed-production") && !strings.Contains(changeItem.From, "deployed-production"){
				  return issueChange.Created.Time, nil
			  }
			}
		}
	}
	return time.Now(), errors.New("No Deploy Time Found...")
}

func laptopToLiveForIssue(issue *IssueHistory) float64 {
  startTime, err := findInProgressTime(issue.Changelog.Histories)
	if err != nil {
		println("No start time found for ", issue.Key)
		return 0.0
	}
	liveTime, err := findDeployTime(issue.Changelog.Histories)
	if err != nil {
		println("No deploy time found for ", issue.Key)
		return 0.0
	}
	duration := liveTime.Sub(startTime)
	return duration.Hours()
}

func CalculateCycleTime(result SearchResult) (float64, float64) {
	var summedCycleTime float64
	var maxCycleTime float64
	summedCycleTime = 0
	maxCycleTime = 0
	issueCount := result.Total
	for _, issue := range result.Issues {
		cycleTime := cycleTimeForIssue(issue)
		if cycleTime > maxCycleTime {
			maxCycleTime = cycleTime
		}
		summedCycleTime += cycleTime
	}
	averageCycleTime := summedCycleTime / float64(issueCount)
	return averageCycleTime, maxCycleTime
}

func CalculateLaptopToLive(changelogs []*IssueHistory) (float64, float64) {
  var summedCycleTime float64
	var maxCycleTime float64
	summedCycleTime = 0
	maxCycleTime = 0
	issueCount := len(changelogs)
	for _, changelog := range changelogs {
		cycleTime := laptopToLiveForIssue(changelog)
		if cycleTime > maxCycleTime {
			maxCycleTime = cycleTime
		}
		summedCycleTime += cycleTime
	}
	averageCycleTime := summedCycleTime / float64(issueCount)
	return averageCycleTime, maxCycleTime
}

func CalculateBugRatio(result SearchResult) float64 {
  return 1.0
}
