package jira

func cycleTimeForIssue(issue Issue) (float64) {
  resolved := issue.Field.Resolved
  duration := resolved.Sub(issue.Field.Created.Time)
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
