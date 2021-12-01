package services

import "golang.design/x/clipboard"

func CopyToClipBoard(s *string) {
	clipboard.Write(clipboard.FmtText, []byte(*s))
}
