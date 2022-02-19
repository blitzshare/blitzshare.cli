package random

import (
	"fmt"
	"strings"

	"github.com/tjarratt/babble"
)

type RndImpl struct {
}

type Rnd interface {
	GenerateRandomWordSequence() *string
}

func New() Rnd {
	return &RndImpl{}
}

func (*RndImpl) GenerateRandomWordSequence() *string {
	babbler := babble.NewBabbler()
	str := fmt.Sprintf("%s-%s", babbler.Babble(), babbler.Babble())
	str = strings.ToLower(strings.Replace(str, "'", "", -1))
	return &str
}
