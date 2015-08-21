package main

import (
  "os"
  "fmt"
  "net/http"
  "net/url"
  "encoding/json"
  "encoding/base64"
  "github.com/codegangsta/cli"
  "github.com/simplereach/timeutils"
)

type IssueFields struct {
  Created timeutils.Time `json:"created"`
  Resolved timeutils.Time `json:"resolutiondate"`
}

type Issue struct {
  Id string `json:"id"`
  Field IssueFields `json:"fields"`
}

type SearchResult struct {
	Total     int    `json:"total"`
	Issues    []Issue `json:"issues"`
}


func cycleTime(c *cli.Context){
  config, cnfErr := LoadConfig()
	if cnfErr != nil {
		println("Configuration Error:", cnfErr)
		return
	}
  jql := fmt.Sprintf("project=%s AND status=Closed AND resolutiondate >= -7d", config.Jira.Project)
  var Url *url.URL
  host := fmt.Sprintf("https://%s.atlassian.net", config.Jira.Subdomain)
  Url, err := url.Parse(host)
  if err != nil {
		println("URL Parsing Error:", err)
		return
	}
  Url.Path += "/rest/api/2/search"
  parameters := url.Values{}
  parameters.Add("jql", jql)
  parameters.Add("maxResults", "500")
  Url.RawQuery = parameters.Encode()
  url := Url.String()

  authString := fmt.Sprintf("%s:%s", config.Jira.Username, config.Jira.Password)
  authHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authString)))
  client := &http.Client{}
  request, err := http.NewRequest("GET", url, nil)
  if err != nil {
		println("Request Building Error:", err)
		return
	}
  request.Header.Add("Content-Type", "application/json")
  request.Header.Add("Authorization", authHeader)

  println("Checking data for last 7 days...")


  response, err := client.Do(request)
  if err != nil {
		println("API Connection Error:", err)
		return
	}

  defer response.Body.Close()
	var result SearchResult
	jsonErr := json.NewDecoder(response.Body).Decode(&result)
  if jsonErr != nil {
		println("API Response Parsing Error:", jsonErr.Error())
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
