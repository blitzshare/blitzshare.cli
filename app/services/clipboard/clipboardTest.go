//go:build test
// +build test

package clipboard

type ClipBoard interface {
	CopyToClipBoard(s *string)
}

type ClipBoardImpl struct {
}

func NewClipBoard() ClipBoard {
	return &ClipBoardImpl{}
}

func (*ClipBoardImpl) CopyToClipBoard(s *string) {
	// Do nothing in test build
}
