package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evizitei/jira-stats/jira"
	"os"
)

func cycleTime(c *cli.Context) {
	config, cnfErr := LoadConfig()
	if cnfErr != nil {
		println("Configuration Error:", cnfErr.Error())
		return
	}

	println("Checking data for last 7 days...")
	jiraClient := jira.Client{Config: config.Jira}
	api := jira.HttpApi{}
	var result jira.SearchResult
	err := jiraClient.QueryRecentlyClosedIssues(api, &result)
	if err != nil {
		return
	}

	average, max := jira.CalculateCycleTime(result)
	println(fmt.Sprintf("Project: %s", config.Jira.Project))
	println(fmt.Sprintf("Username: %s", config.Jira.Username))
	println(fmt.Sprintf("%d total issues resolved", result.Total))
	println(fmt.Sprintf("average cycle time: %f hours", average))
	println(fmt.Sprintf("max cycle time: %f hours", max))
}

func main() {
	app := cli.NewApp()
	app.Name = "jira-stats"
	app.Usage = "Gather metrics about a JIRA project"
	app.Commands = []cli.Command{
		{
			Name:    "cycle-time",
			Aliases: []string{"ct"},
			Usage:   "Calculate how long it takes on average from ticket-entry to done",
			Action:  cycleTime,
		},
	}
	app.Run(os.Args)
}
