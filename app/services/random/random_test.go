package random_test

import (
	"testing"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/random"
	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	words := random.GenerateRandomWords()
	assert.NotNil(t, words)
	assert.Greater(t, len(words), 6)
	assert.NotEqual(t, random.GenerateRandomWords(), words)
}
