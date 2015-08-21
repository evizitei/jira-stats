package main

import (
  "os"
  "fmt"
  "github.com/codegangsta/cli"
  "github.com/evizitei/jira-stats/jira"
)

func cycleTime(c *cli.Context){
  config, cnfErr := LoadConfig()
	if cnfErr != nil {
		println("Configuration Error:", cnfErr.Error())
		return
	}

  println("Checking data for last 7 days...")
  jiraClient := jira.Client{ Config: config.Jira }
  api := jira.HttpApi{}
	var result jira.SearchResult
  err := jiraClient.QueryRecentlyClosedIssues(api, &result)
  if err != nil {
		return
	}

  var summedCycleTime float64
  var maxCycleTime float64
  summedCycleTime = 0
  maxCycleTime = 0
  issueCount := result.Total
  for _, issue := range result.Issues {
    cycleTime := issue.Field.Resolved.Sub(issue.Field.Created.Time).Hours()
    if cycleTime > maxCycleTime {
      maxCycleTime = cycleTime
    }
    summedCycleTime += cycleTime
  }
  averageCycleTime := summedCycleTime / float64(issueCount)

  println(fmt.Sprintf("Project: %s", config.Jira.Project))
  println(fmt.Sprintf("Username: %s", config.Jira.Username))
  println(fmt.Sprintf("%d total issues resolved", result.Total))
  println(fmt.Sprintf("average cycle time: %f hours", averageCycleTime))
  println(fmt.Sprintf("max cycle time: %f hours", maxCycleTime))
}

func main() {
  app := cli.NewApp()
  app.Name = "jira-stats"
  app.Usage = "Gather metrics about a JIRA project"
  app.Commands = []cli.Command{
    {
      Name: "cycle-time",
      Aliases: []string{"ct"},
      Usage: "Calculate how long it takes on average from ticket-entry to done",
      Action: cycleTime,
    },
  }
  app.Run(os.Args)
}
