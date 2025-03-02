package cl

import "fmt"

const (
	Reset        = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
)

// Functions to return colored text
func Red(text string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, text, Reset)
}

func Green(text string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, text, Reset)
}

func Yellow(text string) string {
	return fmt.Sprintf("%s%s%s", ColorYellow, text, Reset)
}

func Blue(text string) string {
	return fmt.Sprintf("%s%s%s", ColorBlue, text, Reset)
}

func Magenta(text string) string {
	return fmt.Sprintf("%s%s%s", ColorMagenta, text, Reset)
}

func Cyan(text string) string {
	return fmt.Sprintf("%s%s%s", ColorCyan, text, Reset)
}

func White(text string) string {
	return fmt.Sprintf("%s%s%s", ColorWhite, text, Reset)
}
