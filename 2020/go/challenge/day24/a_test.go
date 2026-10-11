package day24

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestParsePathValid(t *testing.T) {
	tests := []struct {
		input    string
		expected pos
	}{
		{"e", pos{1, 0}},
		{"w", pos{-1, 0}},
		{"ne", pos{0, 1}},
		{"nw", pos{-1, 1}},
		{"se", pos{1, -1}},
		{"sw", pos{0, -1}},
		{"esew", pos{1, -1}},
		{"nwwswee", pos{0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parsePath(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestParsePathInvalid(t *testing.T) {
	tests := []string{"nq", "sn", "n", "s", "xyz"}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			_, err := parsePath(tt)
			require.Error(t, err)
		})
	}
}

func TestPartA(t *testing.T) {
	example := `sesenwnenenewseeswwswswwnenewsewsw
neeenesenwnwwswnenewnwwsewnenwseswesw
seswneswswsenwwnwse
nwnwneseeswswnenewneswwnewseswneseene
swweswneswnenwsewnwneneseenw
eesenwseswswnenwswnwnwsewwnwsene
sewnenenenesenwsewnenwwwse
wenwwweseeeweswwwnwwe
wsweesenenewnwwnwsenewsenwwsesesenwne
neeswseenwwswnwswswnw
nenwswwsewswnenenewsenwsenwnesesenew
enewnwewneswsewnwswenweswnenwsenwsw
sweneswneswneneenwnewenewwneswswnese
swwesenesewenwneswnwwneseswwne
enesenwswwswneneswsenwnewswseenwsese
wnwnesenesenenwwnenwsewesewsesesew
nenewswnwewswnenesenwnesewesw
eneswnwswnwsenenwnwnwwseeswneewsenese
neswnwewnwnwseenwseesewsenwsweewe
wseweeenwnesenwwwswnew`

	answer, err := A(challenge.FromLiteral(example))
	require.NoError(t, err)
	require.Equal(t, 10, answer)
}

func TestPartAEmpty(t *testing.T) {
	answer, err := A(challenge.FromLiteral(""))
	require.NoError(t, err)
	require.Equal(t, 0, answer)
}

func TestPartASingleTile(t *testing.T) {
	answer, err := A(challenge.FromLiteral("e"))
	require.NoError(t, err)
	require.Equal(t, 1, answer)
}

func TestPartAToggle(t *testing.T) {
	answer, err := A(challenge.FromLiteral("e\ne"))
	require.NoError(t, err)
	require.Equal(t, 0, answer)
}
