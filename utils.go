package gitcomm

import "strings"

func getStoryIdFromLastCommitLog(logMessage string) string {
	lists := strings.Split(logMessage, "\n")
	for _, each := range lists {
		if strings.Contains(each, "--story") {
			parts := strings.Split(each, "=")
			return parts[1]
		}
	}
	return ""
}
