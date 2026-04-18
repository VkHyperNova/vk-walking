package color

import "fmt"

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	Bold   = "\033[1m"
	Italic = "\033[3m"
)

func PrintBoldBlue(text string) string {
	return fmt.Sprintf(Blue + Bold + text + Reset)
}

func PrintBoldYellow(text string) string {
	return fmt.Sprintf(Yellow + Bold + text + Reset)
}
