package output

func Bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func Magenta(s string) string {
	return "\033[35m" + s + "\033[0m"
}

func Faint(s string) string {
	return "\033[2m" + s + "\033[0m"
}
