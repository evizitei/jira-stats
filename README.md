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

4) run your command

`./jira-stats cycle-time`
