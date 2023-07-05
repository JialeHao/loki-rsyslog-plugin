package core

import (
	"errors"
	"regexp"
	"strconv"
)

var severity = [8]string{
	"Emergency",
	"Alert",
	"Critical",
	"Error",
	"Warning",
	"Notice",
	"Informational",
	"Debug",
}

func addLevelTag(m map[string]string, log *[]byte) error {
	var pri int
	var err error

	expr := regexp.MustCompile(`^<(\d+)>.*`)
	matchResult := expr.FindSubmatch(*log)

	if len(matchResult) < 2 {
		err = errors.New("pri match error: can not add level tag")
		return err
	}

	if pri, err = strconv.Atoi(string(matchResult[1])); err != nil {
		return err
	}

	severityCode := pri % 8

	m["level"] = severity[severityCode]
    
	return nil
}
