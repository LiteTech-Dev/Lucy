package output

import "strings"

func bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func mangeta(s string) string {
	return "\033[35m" + s + "\033[0m"
}

func faint(s string) string {
	return "\033[2m" + s + "\033[0m"
}

func captalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
