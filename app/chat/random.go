package chat

import (
	"fmt"

	"github.com/tjarratt/babble"
)

func GenerateRandomWords() string {
	babbler := babble.NewBabbler()
	return fmt.Sprintf("%s-%s", babbler.Babble(), babbler.Babble())
}
