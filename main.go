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

	dateRange := c.String("date-range")
	if dateRange == "default" {
		println("Checking data for last 7 days...")
	} else {
		println("Checking data for", dateRange, "...")
	}

	jiraClient := jira.Client{Config: config.Jira}
	api := jira.HttpApi{}
	var result jira.SearchResult
	err := jiraClient.QueryRecentlyClosedIssues(api, dateRange, &result)
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

	dateRange := c.String("date-range")
	if dateRange == "default" {
		println("Checking data for last 6 weeks...")
	} else {
		println("Checking data for", dateRange, "...")
	}

	jiraClient := jira.Client{Config: config.Jira}
	api := jira.HttpApi{}
	var result jira.SearchResult
	err := jiraClient.QueryRecentlyDeployedIssues(api, dateRange, &result)
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

	dateRange := c.String("date-range")
	if dateRange == "default" {
		println("Checking data for last 7 days...")
	} else {
		println("Checking data for", dateRange, "...")
	}

	jiraClient := jira.Client{Config: config.Jira}
	api := jira.HttpApi{}

	var result jira.SearchResult
	err := jiraClient.QueryRecentlyCreatedIssues(api, dateRange, &result)
	if err != nil {
		return
	}

	bugsOverFeatures := jira.CalculateBugRatio(result)
	println(fmt.Sprintf("Project: %s", config.Jira.Project))
	println(fmt.Sprintf("Username: %s", config.Jira.Username))
	println(fmt.Sprintf("%d total issues created", result.Total))
	println(fmt.Sprintf("bug ratio (bugs/features): %f", bugsOverFeatures))
}

func main() {
	app := cli.NewApp()
	app.Name = "jira-stats"
	app.Usage = "Gather metrics about a JIRA project"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "date-range, d",
			Value: "default",
			Usage: "Date range to check data, the default value is specific to each command but is rational.  Provided like 'YYYY-MM-DD:YYYY-MM-DD' with the earlier date first.",
		},
	}

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
