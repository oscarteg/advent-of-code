use anyhow::{bail, ensure, Context};
use std::collections::VecDeque;

type Part1 = i64;
type Part2 = i64;
type Pos = (i64, i64);
type Vec3 = [i64; 3];

/// Facings are the puzzle's own codes: 0 right, 1 down, 2 left, 3 up.
const DELTAS: [Pos; 4] = [(1, 0), (0, 1), (-1, 0), (0, -1)];

#[derive(Clone, Copy, PartialEq, Eq)]
enum Tile {
    Void,
    Open,
    Wall,
}

#[derive(Clone, Copy)]
enum Move {
    Forward(u32),
    Clockwise,
    CounterClockwise,
}

struct Notes {
    tiles: Vec<Vec<Tile>>,
    width: i64,
    height: i64,
    path: Vec<Move>,
}

impl Notes {
    fn parse(input: &str) -> anyhow::Result<Self> {
        let mut lines = input.lines().map(|line| line.trim_end_matches('\r'));
        let mut tiles = Vec::new();
        for line in lines.by_ref() {
            if line.is_empty() {
                break;
            }
            tiles.push(parse_row(line)?);
        }
        let path = lines
            .find(|line| !line.is_empty())
            .context("notes are missing the path")?;
        let width = tiles.iter().map(Vec::len).max().unwrap_or_default();
        ensure!(width > 0, "notes are missing the map");
        Ok(Self {
            width: i64::try_from(width)?,
            height: i64::try_from(tiles.len())?,
            tiles,
            path: parse_path(path)?,
        })
    }

    fn tile(&self, (x, y): Pos) -> Tile {
        let (Ok(x), Ok(y)) = (usize::try_from(x), usize::try_from(y)) else {
            return Tile::Void;
        };
        self.tiles
            .get(y)
            .and_then(|row| row.get(x))
            .copied()
            .unwrap_or(Tile::Void)
    }

    fn start(&self) -> anyhow::Result<Pos> {
        let x = self
            .tiles
            .first()
            .context("notes are missing the map")?
            .iter()
            .position(|&tile| tile == Tile::Open)
            .context("top row of the map has no open tile")?;
        Ok((i64::try_from(x)?, 0))
    }

    /// Part 1: leave the board and reappear on the opposite edge of the same
    /// row or column, keeping the facing.
    fn wrap_flat(&self, from: Pos, facing: usize) -> (Pos, usize) {
        let (dx, dy) = DELTAS[facing];
        let mut pos = from;
        loop {
            pos = (
                (pos.0 + dx).rem_euclid(self.width),
                (pos.1 + dy).rem_euclid(self.height),
            );
            if self.tile(pos) != Tile::Void {
                return (pos, facing);
            }
        }
    }
}

fn parse_row(line: &str) -> anyhow::Result<Vec<Tile>> {
    line.chars()
        .map(|ch| match ch {
            ' ' => Ok(Tile::Void),
            '.' => Ok(Tile::Open),
            '#' => Ok(Tile::Wall),
            _ => bail!("unexpected map character {ch:?}"),
        })
        .collect()
}

fn parse_path(line: &str) -> anyhow::Result<Vec<Move>> {
    let mut path = Vec::new();
    let mut forward: Option<u32> = None;
    for ch in line.chars() {
        if let Some(digit) = ch.to_digit(10) {
            let count = forward
                .unwrap_or(0)
                .checked_mul(10)
                .and_then(|count| count.checked_add(digit))
                .context("step count is too large")?;
            forward = Some(count);
            continue;
        }
        if let Some(count) = forward.take() {
            path.push(Move::Forward(count));
        }
        path.push(match ch {
            'R' => Move::Clockwise,
            'L' => Move::CounterClockwise,
            _ => bail!("unexpected path character {ch:?}"),
        });
    }
    if let Some(count) = forward {
        path.push(Move::Forward(count));
    }
    ensure!(!path.is_empty(), "path is empty");
    Ok(path)
}

/// One face of the folded cube: where it sits in the net plus the 3D directions
/// its grid axes point in. Every face keeps the invariant normal = down x right,
/// which is what makes the edge crossings below exact.
#[derive(Clone, Copy)]
struct Face {
    origin: Pos,
    right: Vec3,
    down: Vec3,
    normal: Vec3,
}

struct Cube {
    size: i64,
    faces: Vec<Face>,
    by_normal: [usize; 6],
}

impl Cube {
    fn fold(notes: &Notes) -> anyhow::Result<Self> {
        let tiles = notes
            .tiles
            .iter()
            .flatten()
            .filter(|&&tile| tile != Tile::Void)
            .count();
        ensure!(
            tiles > 0 && tiles % 6 == 0,
            "map has {tiles} tiles, which is not six square faces"
        );
        let area = i64::try_from(tiles / 6)?;
        let size = (1..=area)
            .find(|side| side * side >= area)
            .unwrap_or_default();
        ensure!(size * size == area, "{area} tiles per face is not a square");

        let mut cells: Vec<Pos> = Vec::new();
        for y in 0..notes.height {
            for x in 0..notes.width {
                if notes.tile((x, y)) == Tile::Void {
                    continue;
                }
                let cell = (x / size, y / size);
                if !cells.contains(&cell) {
                    cells.push(cell);
                }
            }
        }
        ensure!(
            cells.len() == 6,
            "net has {} faces, expected six",
            cells.len()
        );
        for &(column, row) in &cells {
            for j in 0..size {
                for i in 0..size {
                    ensure!(
                        notes.tile((column * size + i, row * size + j)) != Tile::Void,
                        "face at net cell ({column}, {row}) has holes"
                    );
                }
            }
        }

        let origin_of = |(column, row): Pos| (column * size, row * size);
        let mut assigned: Vec<Option<Face>> = vec![None; cells.len()];
        assigned[0] = Some(Face {
            origin: origin_of(cells[0]),
            right: [1, 0, 0],
            down: [0, 1, 0],
            normal: [0, 0, -1],
        });
        let mut queue = VecDeque::from([0usize]);
        let mut folded = 1;
        while let Some(index) = queue.pop_front() {
            let Some(face) = assigned[index] else {
                continue;
            };
            for (facing, &(dx, dy)) in DELTAS.iter().enumerate() {
                let neighbor = (cells[index].0 + dx, cells[index].1 + dy);
                let Some(next) = cells.iter().position(|&cell| cell == neighbor) else {
                    continue;
                };
                if assigned[next].is_some() {
                    continue;
                }
                assigned[next] = Some(face.fold_across(facing, origin_of(cells[next])));
                folded += 1;
                queue.push_back(next);
            }
        }
        ensure!(
            folded == 6,
            "net does not fold into a cube: it is not connected"
        );

        let faces: Vec<Face> = assigned.into_iter().flatten().collect();
        let mut by_normal = [usize::MAX; 6];
        for (index, face) in faces.iter().enumerate() {
            let slot = &mut by_normal[axis_index(face.normal)];
            ensure!(
                *slot == usize::MAX,
                "net does not fold into a cube: two faces share a side"
            );
            *slot = index;
        }
        Ok(Self {
            size,
            faces,
            by_normal,
        })
    }

    fn face_at(&self, (x, y): Pos) -> Option<Face> {
        self.faces.iter().copied().find(|face| {
            (face.origin.0..face.origin.0 + self.size).contains(&x)
                && (face.origin.1..face.origin.1 + self.size).contains(&y)
        })
    }

    /// Part 2: step over a cube edge. Cell (i, j) of a face sits at the doubled
    /// lattice point size * normal + (2i - size + 1) * right + (2j - size + 1) *
    /// down. Crossing an edge swaps the components along the old normal and the
    /// travel direction, which reduces to adding travel - normal, and the walk
    /// continues along -normal on the face whose normal is the travel direction.
    fn wrap(&self, from: Pos, facing: usize) -> anyhow::Result<(Pos, usize)> {
        let face = self.face_at(from).context("walked off the cube")?;
        let travel = match facing {
            0 => face.right,
            1 => face.down,
            2 => negate(face.right),
            _ => negate(face.down),
        };
        let (i, j) = (from.0 - face.origin.0, from.1 - face.origin.1);
        let here = add(
            scale(face.normal, self.size),
            add(
                scale(face.right, 2 * i - self.size + 1),
                scale(face.down, 2 * j - self.size + 1),
            ),
        );
        let there = add(here, subtract(travel, face.normal));

        let dest = self.faces[self.by_normal[axis_index(travel)]];
        let local = subtract(there, scale(dest.normal, self.size));
        let column = (dot(local, dest.right) + self.size - 1) / 2;
        let row = (dot(local, dest.down) + self.size - 1) / 2;
        let entry = negate(face.normal);
        let entry_facing = if entry == dest.right {
            0
        } else if entry == dest.down {
            1
        } else if entry == negate(dest.right) {
            2
        } else {
            3
        };
        Ok(((dest.origin.0 + column, dest.origin.1 + row), entry_facing))
    }
}

impl Face {
    /// The net neighbour in the given facing rotates this basis a quarter turn
    /// around the shared edge.
    fn fold_across(&self, facing: usize, origin: Pos) -> Self {
        let (right, down, normal) = match facing {
            0 => (negate(self.normal), self.down, self.right),
            1 => (self.right, negate(self.normal), self.down),
            2 => (self.normal, self.down, negate(self.right)),
            _ => (self.right, self.normal, negate(self.down)),
        };
        Self {
            origin,
            right,
            down,
            normal,
        }
    }
}

fn add(a: Vec3, b: Vec3) -> Vec3 {
    [a[0] + b[0], a[1] + b[1], a[2] + b[2]]
}

fn subtract(a: Vec3, b: Vec3) -> Vec3 {
    [a[0] - b[0], a[1] - b[1], a[2] - b[2]]
}

fn negate(a: Vec3) -> Vec3 {
    [-a[0], -a[1], -a[2]]
}

fn scale(a: Vec3, factor: i64) -> Vec3 {
    [a[0] * factor, a[1] * factor, a[2] * factor]
}

fn dot(a: Vec3, b: Vec3) -> i64 {
    a[0] * b[0] + a[1] * b[1] + a[2] * b[2]
}

/// Slot for one of the six axis-aligned unit normals.
fn axis_index(normal: Vec3) -> usize {
    let axis = normal
        .iter()
        .position(|&part| part != 0)
        .unwrap_or_default();
    axis * 2 + usize::from(normal[axis] < 0)
}

fn walk<W>(notes: &Notes, mut wrap: W) -> anyhow::Result<i64>
where
    W: FnMut(Pos, usize) -> anyhow::Result<(Pos, usize)>,
{
    let mut pos = notes.start()?;
    let mut facing = 0;
    for step in &notes.path {
        match *step {
            Move::Clockwise => facing = (facing + 1) % 4,
            Move::CounterClockwise => facing = (facing + 3) % 4,
            Move::Forward(count) => {
                for _ in 0..count {
                    let (dx, dy) = DELTAS[facing];
                    let ahead = (pos.0 + dx, pos.1 + dy);
                    let (next, next_facing) = if notes.tile(ahead) == Tile::Void {
                        wrap(pos, facing)?
                    } else {
                        (ahead, facing)
                    };
                    if notes.tile(next) == Tile::Wall {
                        break;
                    }
                    (pos, facing) = (next, next_facing);
                }
            }
        }
    }
    Ok(1000 * (pos.1 + 1) + 4 * (pos.0 + 1) + i64::try_from(facing)?)
}

pub fn solve(input: &str) -> anyhow::Result<(Part1, Part2)> {
    let notes = Notes::parse(input)?;
    let part1 = walk(&notes, |pos, facing| Ok(notes.wrap_flat(pos, facing)))?;
    let cube = Cube::fold(&notes)?;
    let part2 = walk(&notes, |pos, facing| cube.wrap(pos, facing))?;
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day22-example.txt");

    /// A 1-4-1 net with one tile per face: a layout nothing like the example's.
    const UNIT_NET: &str = " .\n....\n  .\n\n2";
    const UNIT_NET_WALL: &str = " .\n....\n  #\n\n2";

    #[test]
    fn part1_example() {
        let (part1, _) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 6032);
    }

    #[test]
    fn part2_example() {
        let (_, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part2, 5031);
    }

    #[test]
    fn folds_a_net_with_single_tile_faces() {
        // Flat wrapping never leaves the top row; on the cube the walk crosses
        // onto the right face and then onto the front face.
        assert_eq!(solve(UNIT_NET).unwrap(), (1008, 3013));
    }

    #[test]
    fn wall_stops_the_walk_after_a_crossing() {
        assert_eq!(solve(UNIT_NET_WALL).unwrap(), (1008, 2013));
    }

    #[test]
    fn rejects_missing_path() {
        assert!(solve("....\n....").is_err());
    }

    #[test]
    fn rejects_unknown_characters() {
        assert!(solve(" .\n..x.\n  .\n\n2").is_err());
        assert!(solve(" .\n....\n  .\n\n2X").is_err());
    }

    #[test]
    fn rejects_maps_that_are_not_cube_nets() {
        assert!(solve("....\n....\n\n1").is_err());
        assert!(solve("......\n......\n\n1").is_err());
    }

    #[test]
    fn rejects_empty_input() {
        assert!(solve("").is_err());
    }
}
