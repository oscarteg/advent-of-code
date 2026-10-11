package day20

import (
	"slices"
	"strings"
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = `Tile 2311:
..##.#..#.
##..#.....
#...##..#.
####.#...#
##.##.###.
##...#.###
.#.#.#..##
..#....#..
###...#.#.
..###..###

Tile 1951:
#.##...##.
#.####...#
.....#..##
#...######
.##.#....#
.###.#####
###.##.##.
.###....#.
..#.#..#.#
#...##.#..

Tile 1171:
####...##.
#..##.#..#
##.#..#.#.
.###.####.
..###.####
.##....##.
.#...####.
#.##.####.
####..#...
.....##...

Tile 1427:
###.##.#..
.#..#.##..
.#.##.#..#
#.#.#.##.#
....#...##
...##..##.
...#.#####
.#.####.#.
..#..###.#
..##.#..#.

Tile 1489:
##.#.#....
..##...#..
.##..##...
..#...#...
#####...#.
#..#.#.#.#
...#.#.#..
##.#...##.
..##.##.##
###.##.#..

Tile 2473:
#....####.
#..#.##...
#.##..#...
######.#.#
.#...#.#.#
.#########
.###.#..#.
########.#
##...##.#.
..###.#.#.

Tile 2971:
..#.#....#
#...###...
#.#.###...
##.##..#..
.#####..##
.#..####.#
#..#.#..#.
..####.###
..#.#.###.
...#.#.#.#

Tile 2729:
...#.#.#.#
####.#....
..#.#.....
....#..#.#
.##..##.#.
.#.####...
####.#.#..
##.####...
##..#.##..
#.##...##.

Tile 3079:
#.#.#####.
.#..######
..#.......
######....
####.#..#.
.#...#.##.
#.#####.##
..#.###...
..#.......
..#.###...
`

// flipTiles mirrors every tile; the tiles arrive in a random orientation anyway.
func flipTiles(input string) string {
	lines := strings.Split(input, "\n")
	for index, line := range lines {
		if strings.HasPrefix(line, "Tile ") {
			continue
		}
		mirrored := []byte(line)
		for left, right := 0, len(mirrored)-1; left < right; left, right = left+1, right-1 {
			mirrored[left], mirrored[right] = mirrored[right], mirrored[left]
		}
		lines[index] = string(mirrored)
	}
	return strings.Join(lines, "\n")
}

func TestA(t *testing.T) {
	answer, err := A(challenge.FromLiteral(example))

	require.NoError(t, err)
	require.Equal(t, 20899048083289, answer)
}

func TestAIgnoresTileOrderAndOrientation(t *testing.T) {
	blocks := strings.Split(strings.TrimSpace(flipTiles(example)), "\n\n")
	slices.Reverse(blocks)

	answer, err := A(challenge.FromLiteral(strings.Join(blocks, "\n\n")))

	require.NoError(t, err)
	require.Equal(t, 20899048083289, answer)
}

func TestAFourRotationsReturnTheOriginalImage(t *testing.T) {
	img, err := parseImage([]string{"#..", ".#.", "..#"})
	require.NoError(t, err)

	require.Equal(t, img, img.rotated().rotated().rotated().rotated())
	require.Equal(t, edge("..#"), img.rotated().top())
	require.Equal(t, edge("..#"), img.flipped().top())
	require.Equal(t, [][]bool{{true}}, img.interior().rows)
}

func TestANoTiles(t *testing.T) {
	_, err := A(challenge.FromLiteral("   \n\n"))

	require.Error(t, err)
}

func TestAMissingTileHeader(t *testing.T) {
	_, err := A(challenge.FromLiteral("..##.#..#.\n##..#.....\n#...##..#.\n"))

	require.Error(t, err)
}

func TestAInvalidTileID(t *testing.T) {
	_, err := A(challenge.FromLiteral(strings.Replace(example, "Tile 2311:", "Tile two:", 1)))

	require.Error(t, err)
}

func TestAIncompleteImage(t *testing.T) {
	blocks := strings.Split(strings.TrimSpace(example), "\n\n")

	_, err := A(challenge.FromLiteral(strings.Join(blocks[:8], "\n\n")))

	require.Error(t, err)
}
