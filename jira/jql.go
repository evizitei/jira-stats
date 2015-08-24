package jira

import (
	"fmt"
	"strings"
)

// Jql is a simple object for building up a JQL string before shipping it off
//  to the JIRA API.  The date range "default" will let each operation choose
//  it's own default range, but "YYYY-MM-DD:YYYY-MM-DD" will give an explicit
//  date range for any constructed operation to work over.
type Jql struct {
	Project   string
	DateRange string
}

func (j Jql) jqlForProject(whereClause string) string {
	baseJql := "project=%s AND %s"
	return fmt.Sprintf(baseJql, j.Project, whereClause)
}

func (j Jql) buildDateClause(field string, defaultValue string) string {
	if j.DateRange != "default" {
		dates := strings.Split(j.DateRange, ":")
		return fmt.Sprintf("%s >= %s AND %s <= %s", field, dates[0], field, dates[1])
	}
	return fmt.Sprintf("%s >= %s", field, defaultValue)
}

// Closed builds a JQL string with a clause checking for closing events within
//  the date range
func (j Jql) Closed() string {
	dateClause := j.buildDateClause("resolutiondate", "-7d")
	whereClause := fmt.Sprintf("status=Closed AND %s", dateClause)
	return j.jqlForProject(whereClause)
}

// Deployed builds a JQL string with a clause checking for the deployed-production label
//  on tickets resovled within the date range
func (j Jql) Deployed() string {
	dateClause := j.buildDateClause("resolutiondate", "-42d")
	whereClause := fmt.Sprintf("labels in (deployed-production) AND %s", dateClause)
	return j.jqlForProject(whereClause)
}

// Created builds a JQL string with a clause checking for creation events within
//  the date range
func (j Jql) Created() string {
	dateClause := j.buildDateClause("createdDate", "-7d")
	return j.jqlForProject(dateClause)
}
