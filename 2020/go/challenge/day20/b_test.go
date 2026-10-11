package day20

import (
	"strings"
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	answer, err := B(challenge.FromLiteral(example))

	require.NoError(t, err)
	require.Equal(t, 273, answer)
}

func TestBIgnoresTileOrientation(t *testing.T) {
	answer, err := B(challenge.FromLiteral(flipTiles(example)))

	require.NoError(t, err)
	require.Equal(t, 273, answer)
}

func TestBRaggedTile(t *testing.T) {
	ragged := strings.Replace(example, "..##.#..#.\n", "..##.#..#\n", 1)

	_, err := B(challenge.FromLiteral(ragged))

	require.Error(t, err)
}

func TestBMalformedInput(t *testing.T) {
	_, err := B(challenge.FromLiteral("Tile 1:\n..x\n...\n..."))

	require.Error(t, err)
}
