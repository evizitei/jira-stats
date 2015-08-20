package main

import (
  "os"
  "github.com/codegangsta/cli"
)

func cycleTime(c *cli.Context){
  config, cnfErr := LoadConfig()
	if cnfErr != nil {
		println("Configuration Error:", cnfErr)
		return
	}
  println("Contacting", config.Jira.Subdomain, "on behalf of", config.Jira.Username)
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
