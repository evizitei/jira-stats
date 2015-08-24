package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/evizitei/jira-stats/jira"
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

func laptopToLive(c *cli.Context) {
	config, cnfErr := LoadConfig()
	if cnfErr != nil {
		println("Configuration Error:", cnfErr.Error())
		return
	}

	println("Checking data for last 6 weeks...")
	jiraClient := jira.Client{Config: config.Jira}
	api := jira.HttpApi{}
	var result jira.SearchResult
	err := jiraClient.QueryRecentlyDeployedIssues(api, &result)
	changelogs, err := jiraClient.QueryChangelogsForResultSet(api, &result)
	average, max := jira.CalculateLaptopToLive(changelogs)
	if err != nil {
		return
	}
	println(fmt.Sprintf("Project: %s", config.Jira.Project))
	println(fmt.Sprintf("Username: %s", config.Jira.Username))
	println(fmt.Sprintf("%d total issues deployed", result.Total))
	println(fmt.Sprintf("average laptop-to-live time: %f hours", average))
	println(fmt.Sprintf("max laptop-to-live time: %f hours", max))
}

func bugRatio(c *cli.Context) {
	config, cnfErr := LoadConfig()
	if cnfErr != nil {
		println("Configuration Error:", cnfErr.Error())
		return
	}

	println("Checking data for last 7 days...")
	jiraClient := jira.Client{Config: config.Jira}
	api := jira.HttpApi{}

	var result jira.SearchResult
	err := jiraClient.QueryRecentlyCreatedIssues(api, &result)
	if err != nil {
		return
	}

	bugsOverFeatures := jira.CalculateBugRatio(result)
	println(fmt.Sprintf("Project: %s", config.Jira.Project))
	println(fmt.Sprintf("Username: %s", config.Jira.Username))
	println(fmt.Sprintf("%d total issues resolved", result.Total))
	println(fmt.Sprintf("bug ratio (bugs/features): %f", bugsOverFeatures))
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
		{
			Name:    "laptop-to-live",
			Aliases: []string{"ltl"},
			Usage:   "Calculate cycle time for 'someone started' to 'on production'",
			Action:  laptopToLive,
		},
		{
			Name:    "bug-ratio",
			Aliases: []string{"br"},
			Usage:   "give the ratio of bugs to feature tickets created",
			Action:  bugRatio,
		},
	}
	app.Run(os.Args)
}
