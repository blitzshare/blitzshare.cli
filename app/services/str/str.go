package str

import "strings"

func SanatizeStr(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}
