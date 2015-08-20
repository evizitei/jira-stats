package main

import (
  "os"
  "fmt"
  "net/http"
  "net/url"
  "encoding/base64"
  "io/ioutil"
  "github.com/codegangsta/cli"
)

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
  Url.RawQuery = parameters.Encode()
  url := Url.String()
  println("hitting this endpoint", url)

  authString := fmt.Sprintf("%s:%s", config.Jira.Username, config.Jira.Password)
  authHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authString)))
  println("Contacting", url, "on behalf of", config.Jira.Username)
  client := &http.Client{}
  request, err := http.NewRequest("GET", url, nil)
  if err != nil {
		println("Request Building Error:", err)
		return
	}
  request.Header.Add("Content-Type", "application/json")
  request.Header.Add("Authorization", authHeader)
  response, err := client.Do(request)
  if err != nil {
		println("API Connection Error:", err)
		return
	}

  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
		println("API Response Parsing Error:", err)
		return
	}
  println("We have a response: ", fmt.Sprintf("%s\n", string(contents)))
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
