package functions

import "strings"

func ParseFileKey(DownloadPath string) string {
	splits := strings.Split(DownloadPath, "/")
	return splits[len(splits)-1]
}
