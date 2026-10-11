package challenge

type Grid struct {
	tiles         []string
	width, height int
}

// This is not really a "grid" because it is not a multi dimensional array but
// I don't know what consumes more memory. An array of string or a multi dimensional array of characters.
func NewGrid(width, height int) *Grid {
	return &Grid{
		tiles:  make([]string, width*height),
		height: height,
		width:  width,
	}
}

func (grid *Grid) CountTreesInPath(right, down int) int {
	x := 0
	var lineLength int
	countTrees := 0
	for y := 0; y < len(grid.tiles); y += down {
		if y != 0 {
			x += right
			if string(grid.tiles[y][x%lineLength]) == "#" {
				countTrees++
			}
		} else {
			lineLength = len(grid.tiles[y])
			continue
		}

	}
	return countTrees
}
