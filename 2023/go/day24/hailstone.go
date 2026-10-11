package main

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var errInvalidHailstones = errors.New("invalid hailstones")

// The puzzle's test area for part 1, inclusive on both axes.
const (
	testAreaLow  = 200_000_000_000_000
	testAreaHigh = 400_000_000_000_000
)

// vector holds exact rational coordinates, because both the crossings of part 1
// and the rock of part 2 need more precision than float64 or int64 provide.
type vector struct {
	x, y, z *big.Rat
}

type hailstone struct {
	position, velocity vector
}

func (v vector) axis(index int) *big.Rat {
	switch index {
	case 0:
		return v.x
	case 1:
		return v.y
	default:
		return v.z
	}
}

func (v vector) sub(other vector) vector {
	return vector{
		x: new(big.Rat).Sub(v.x, other.x),
		y: new(big.Rat).Sub(v.y, other.y),
		z: new(big.Rat).Sub(v.z, other.z),
	}
}

func (v vector) cross(other vector) vector {
	return vector{
		x: determinant(v.y, v.z, other.y, other.z),
		y: determinant(v.z, v.x, other.z, other.x),
		z: determinant(v.x, v.y, other.x, other.y),
	}
}

func (v vector) isZero() bool {
	return v.x.Sign() == 0 && v.y.Sign() == 0 && v.z.Sign() == 0
}

func (v vector) coordinateSum() *big.Rat {
	return new(big.Rat).Add(v.x, new(big.Rat).Add(v.y, v.z))
}

// determinant returns a*d - b*c, the determinant of [[a, b], [c, d]].
func determinant(a, b, c, d *big.Rat) *big.Rat {
	return new(big.Rat).Sub(new(big.Rat).Mul(a, d), new(big.Rat).Mul(b, c))
}

func parseHailstones(input string) ([]hailstone, error) {
	var stones []hailstone
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		stone, err := parseHailstone(line)
		if err != nil {
			return nil, err
		}
		stones = append(stones, stone)
	}
	if len(stones) == 0 {
		return nil, fmt.Errorf("%w: input holds no hailstones", errInvalidHailstones)
	}
	return stones, nil
}

func parseHailstone(line string) (hailstone, error) {
	position, velocity, found := strings.Cut(line, "@")
	if !found {
		return hailstone{}, fmt.Errorf("%w: %q needs an '@' separator", errInvalidHailstones, line)
	}
	parsedPosition, err := parseVector(position)
	if err != nil {
		return hailstone{}, fmt.Errorf("position: %w", err)
	}
	parsedVelocity, err := parseVector(velocity)
	if err != nil {
		return hailstone{}, fmt.Errorf("velocity: %w", err)
	}
	return hailstone{position: parsedPosition, velocity: parsedVelocity}, nil
}

func parseVector(input string) (vector, error) {
	fields := strings.Split(input, ",")
	if len(fields) != 3 {
		return vector{}, fmt.Errorf("%w: %q needs three coordinates", errInvalidHailstones, strings.TrimSpace(input))
	}
	var coordinates [3]*big.Rat
	for index, field := range fields {
		field = strings.TrimSpace(field)
		value, ok := new(big.Int).SetString(field, 10)
		if !ok {
			return vector{}, fmt.Errorf("%w: %q is not an integer", errInvalidHailstones, field)
		}
		coordinates[index] = new(big.Rat).SetInt(value)
	}
	return vector{x: coordinates[0], y: coordinates[1], z: coordinates[2]}, nil
}

// countCrossings counts the pairs of hailstones whose forward paths cross inside
// the square [low, high] on both axes, ignoring the z axis.
func countCrossings(stones []hailstone, low, high int64) int64 {
	lowBound, highBound := big.NewRat(low, 1), big.NewRat(high, 1)
	var crossings int64
	for index, first := range stones {
		for _, second := range stones[index+1:] {
			if crossesInArea(first, second, lowBound, highBound) {
				crossings++
			}
		}
	}
	return crossings
}

func crossesInArea(first, second hailstone, low, high *big.Rat) bool {
	// Cramer's rule on first.position + s*first.velocity == second.position +
	// t*second.velocity, restricted to the x and y axes.
	denominator := determinant(first.velocity.x, first.velocity.y, second.velocity.x, second.velocity.y)
	if denominator.Sign() == 0 {
		return false
	}
	offset := second.position.sub(first.position)
	firstTime := new(big.Rat).Quo(determinant(offset.x, offset.y, second.velocity.x, second.velocity.y), denominator)
	secondTime := new(big.Rat).Quo(determinant(offset.x, offset.y, first.velocity.x, first.velocity.y), denominator)
	if firstTime.Sign() <= 0 || secondTime.Sign() <= 0 {
		return false
	}
	for _, path := range [][2]*big.Rat{
		{first.position.x, first.velocity.x},
		{first.position.y, first.velocity.y},
	} {
		coordinate := new(big.Rat).Add(path[0], new(big.Rat).Mul(firstTime, path[1]))
		if coordinate.Cmp(low) < 0 || coordinate.Cmp(high) > 0 {
			return false
		}
	}
	return true
}

// relativePath is a hailstone's path in the first hailstone's frame of
// reference, where the rock's own path runs through the origin. The normal is
// the one of the plane spanned by this path and that origin.
type relativePath struct {
	offset, velocity, normal vector
}

func (p relativePath) positionAt(time *big.Rat, axis int) *big.Rat {
	return new(big.Rat).Add(p.offset.axis(axis), new(big.Rat).Mul(time, p.velocity.axis(axis)))
}

// rockPositionSum sums the coordinates of the position the rock has to be thrown
// from. Seen from the first hailstone the rock's path runs through the origin,
// so every other path spans a plane with it. The rock's direction lies in all of
// those planes, which makes it the cross product of two of their normals.
// Crossing that line with two paths gives two hit times and two points: their
// difference over the time span is the rock's velocity, and walking back to time
// zero gives its start position.
func rockPositionSum(stones []hailstone) (int64, error) {
	if len(stones) == 0 {
		return 0, fmt.Errorf("%w: input holds no hailstones", errInvalidHailstones)
	}
	base := stones[0]
	first, second, direction, err := rockDirection(relativePaths(base, stones[1:]))
	if err != nil {
		return 0, err
	}
	firstTime, err := hitTime(first, direction)
	if err != nil {
		return 0, err
	}
	secondTime, err := hitTime(second, direction)
	if err != nil {
		return 0, err
	}
	span := new(big.Rat).Sub(secondTime, firstTime)
	if span.Sign() == 0 {
		return 0, fmt.Errorf("%w: two hailstones are hit at the same time", errInvalidHailstones)
	}
	total := base.position.coordinateSum()
	for axis := 0; axis < 3; axis++ {
		firstPoint := first.positionAt(firstTime, axis)
		secondPoint := second.positionAt(secondTime, axis)
		velocity := new(big.Rat).Quo(new(big.Rat).Sub(secondPoint, firstPoint), span)
		total.Add(total, new(big.Rat).Sub(firstPoint, new(big.Rat).Mul(firstTime, velocity)))
	}
	if !total.IsInt() {
		return 0, fmt.Errorf("%w: the rock does not start on integer coordinates", errInvalidHailstones)
	}
	if !total.Num().IsInt64() {
		return 0, fmt.Errorf("%w: the rock position sum does not fit in an int64", errInvalidHailstones)
	}
	return total.Num().Int64(), nil
}

func relativePaths(base hailstone, rest []hailstone) []relativePath {
	paths := make([]relativePath, 0, len(rest))
	for _, stone := range rest {
		offset := stone.position.sub(base.position)
		velocity := stone.velocity.sub(base.velocity)
		normal := offset.cross(velocity)
		// A path through the origin spans no plane with it, so it says nothing
		// about the rock's direction.
		if normal.isZero() {
			continue
		}
		paths = append(paths, relativePath{offset: offset, velocity: velocity, normal: normal})
	}
	return paths
}

// rockDirection picks two paths whose planes meet in the rock's line and returns
// that line's direction alongside them.
func rockDirection(paths []relativePath) (relativePath, relativePath, vector, error) {
	for index, first := range paths {
		for _, second := range paths[index+1:] {
			if direction := first.normal.cross(second.normal); !direction.isZero() {
				return first, second, direction, nil
			}
		}
	}
	return relativePath{}, relativePath{}, vector{},
		fmt.Errorf("%w: the hailstones do not determine a single rock", errInvalidHailstones)
}

// hitTime returns when the path crosses the rock's line: the point
// offset + time*velocity lies on that line exactly when its cross product with
// the direction is zero.
func hitTime(path relativePath, direction vector) (*big.Rat, error) {
	offset := path.offset.cross(direction)
	velocity := path.velocity.cross(direction)
	for axis := 0; axis < 3; axis++ {
		if velocity.axis(axis).Sign() != 0 {
			return new(big.Rat).Quo(new(big.Rat).Neg(offset.axis(axis)), velocity.axis(axis)), nil
		}
	}
	return nil, fmt.Errorf("%w: a hailstone runs parallel to the rock", errInvalidHailstones)
}

func solve(input string) (part1, part2 int64, err error) {
	return solveWith(input, testAreaLow, testAreaHigh)
}

// solveWith counts the part 1 crossings inside an arbitrary test area, so the
// official example can use its own small square.
func solveWith(input string, low, high int64) (part1, part2 int64, err error) {
	stones, err := parseHailstones(input)
	if err != nil {
		return 0, 0, err
	}
	rock, err := rockPositionSum(stones)
	if err != nil {
		return 0, 0, err
	}
	return countCrossings(stones, low, high), rock, nil
}
