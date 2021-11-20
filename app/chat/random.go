package chat

import (
	"fmt"
	"strings"

	"github.com/tjarratt/babble"
)

func GenerateRandomWords() string {
	babbler := babble.NewBabbler()
	str := fmt.Sprintf("%s-%s", babbler.Babble(), babbler.Babble())
	return strings.ToLower(str)
}
