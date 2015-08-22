jira-stats
===============

An analysis tool for getting some interesting metrics around your jira project

## Usage

1) pull down the repo:

`go get github.com/evizitei/jira-stats`

2) go to your repo

`cd $GOPATH/src/github.com/evizitei/jira-stats`

2) build your config file

`cp jirastats.gcfg.example jirastats.gcfg`

`vim jirastats.gcfg`

3) build your binary

`go build`

4) see what commands are available

`./jira-stats`

5) run your command

`./jira-stats cycle-time`

## Running Tests

once you've pulled down the repo, package tests for the jira client package
can be run with:

`go test github.com/evizitei/jira-stats/jira`
