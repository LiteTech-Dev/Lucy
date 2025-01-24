package output

import "strings"

func bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func capitalize(v interface{}) string {
	s := v.(string)
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func faint(s string) string {
	return "\033[2m" + s + "\033[0m"
}

func red(s string) string {
	return "\033[31m" + s + "\033[0m"
}

func green(s string) string {
	return "\033[32m" + s + "\033[0m"
}

func yellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}

func blue(s string) string {
	return "\033[34m" + s + "\033[0m"
}

func mangeta(s string) string {
	return "\033[35m" + s + "\033[0m"
}

func cyan(s string) string {
	return "\033[36m" + s + "\033[0m"
}
