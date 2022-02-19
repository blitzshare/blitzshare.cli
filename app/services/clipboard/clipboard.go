//go:build !test
// +build !test

package clipboard

import "golang.design/x/clipboard"

type ClipBoard interface {
	CopyToClipBoard(s *string)
}

type ClipBoardImpl struct {
}

func New() ClipBoard {
	return &ClipBoardImpl{}
}

func (*ClipBoardImpl) CopyToClipBoard(s *string) {
	clipboard.Write(clipboard.FmtText, []byte(*s))
}
