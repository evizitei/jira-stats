package jira

import (
	"errors"
	"strings"
	"time"
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
			if changeItem.Field == "labels" {
				if strings.Contains(changeItem.To, "deployed-production") && !strings.Contains(changeItem.From, "deployed-production") {
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
		return 0.0
	}
	liveTime, err := findDeployTime(issue.Changelog.Histories)
	if err != nil {
		return 0.0
	}
	duration := liveTime.Sub(startTime)
	return duration.Hours()
}

func averageAndMax(values []float64) (average float64, max float64) {
	var sum float64
	count := len(values)
	for _, val := range values {
		if val == 0.0 {
			count--
		}
		if val > max {
			max = val
		}
		sum += val
	}
	average = sum / float64(count)
	return average, max
}

func bugRatio(issueTypes []string) float64 {
	bugCount := 0
	featureCount := 0
	for _, issueType := range issueTypes {
		if issueType == "Bug" {
			bugCount++
		} else if issueType == "New Feature" {
			featureCount++
		}
	}
	if featureCount == 0 || bugCount == 0 {
		println("Insufficient data to return a meaningful ratio")
		return 0.0
	}
	return float64(bugCount) / float64(featureCount)
}

// CalculateCycleTime iterates over the issues in a search result
//  from a JIRA api query and checks the created time and resolved time
//  for each one, giving an average for how long it takes from getting a ticket
//  filed to having the work on that ticket complete.  Lower numbers are better.
func CalculateCycleTime(result SearchResult) (float64, float64) {
	cycleTimes := make([]float64, len(result.Issues))
	for i, issue := range result.Issues {
		cycleTimes[i] = cycleTimeForIssue(issue)
	}
	return averageAndMax(cycleTimes)
}

// CalculateLaptopToLive accepts an array of changelogs, one for each
//  issue in a result set, and checks the label application times and
//  start times of each one to figure out how long each issue takes to
//  go from "work started" to "deployed to production".  The result is in numbers,
//  lower numbers are better.
func CalculateLaptopToLive(changelogs []*IssueHistory) (float64, float64) {
	ltlTimes := make([]float64, len(changelogs))
	for i, changelog := range changelogs {
		ltlTimes[i] = laptopToLiveForIssue(changelog)
	}
	return averageAndMax(ltlTimes)
}

// CalculateBugRatio takes a SearchResult from a jira API query and
//   counts the bugs and features in that set, then divides bugs
//	 by features to get a ratio.  0.0 would be amazing.  < 1.0 means more features
//	 then bugs.  > 1.0 means more bugs than features.
func CalculateBugRatio(result SearchResult) float64 {
	issueTypes := make([]string, len(result.Issues))
	for i, issue := range result.Issues {
		issueTypes[i] = issue.Field.IssueType.Name
	}
	return bugRatio(issueTypes)
}
