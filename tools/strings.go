package tools

import "strings"

var cutset = " \n\t"

func TrimBlankChar(str string) string {
	return strings.Trim(str, cutset)
}
