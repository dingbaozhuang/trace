package util

import "regexp"

// greedy match JSON from string
func GreedyMatchJSONString(s string) string {
	reg := regexp.MustCompile(`{.*}`)
	matchs := reg.FindStringSubmatch(s)
	if len(matchs) < 1 {
		return ""
	}

	return matchs[0]
}
