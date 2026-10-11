package day24

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestPartB(t *testing.T) {
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

	answer, err := B(challenge.FromLiteral(example))
	require.NoError(t, err)
	require.Equal(t, 2208, answer)
}

func TestPartBEmpty(t *testing.T) {
	answer, err := B(challenge.FromLiteral(""))
	require.NoError(t, err)
	require.Equal(t, 0, answer)
}

func TestPartBSingleTile(t *testing.T) {
	answer, err := B(challenge.FromLiteral("e"))
	require.NoError(t, err)
	require.GreaterOrEqual(t, answer, 0)
}
