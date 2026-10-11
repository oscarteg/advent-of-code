use crate::day10::custom_error::AocError;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub(crate) enum Direction {
    North,
    South,
    East,
    West,
}

impl Direction {
    fn opposite(self) -> Direction {
        match self {
            Direction::North => Direction::South,
            Direction::South => Direction::North,
            Direction::East => Direction::West,
            Direction::West => Direction::East,
        }
    }

    fn delta(self) -> (isize, isize) {
        match self {
            Direction::North => (-1, 0),
            Direction::South => (1, 0),
            Direction::East => (0, 1),
            Direction::West => (0, -1),
        }
    }
}

const ALL_DIRECTIONS: [Direction; 4] = [
    Direction::North,
    Direction::South,
    Direction::East,
    Direction::West,
];

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum Tile {
    Pipe(Direction, Direction),
    Ground,
    Start,
}

impl Tile {
    fn from_char(c: char) -> Result<Tile, AocError> {
        use Direction::*;
        Ok(match c {
            '|' => Tile::Pipe(North, South),
            '-' => Tile::Pipe(East, West),
            'L' => Tile::Pipe(North, East),
            'J' => Tile::Pipe(North, West),
            '7' => Tile::Pipe(South, West),
            'F' => Tile::Pipe(South, East),
            '.' => Tile::Ground,
            'S' => Tile::Start,
            other => return Err(AocError::InvalidTile(other)),
        })
    }

    fn connects(self, direction: Direction) -> bool {
        matches!(self, Tile::Pipe(a, b) if a == direction || b == direction)
    }
}

/// A rectangular sketch of pipe tiles, as described by the puzzle input.
pub(crate) struct Grid {
    tiles: Vec<Vec<Tile>>,
}

impl Grid {
    pub(crate) fn parse(input: &str) -> Result<Grid, AocError> {
        let lines: Vec<&str> = input
            .lines()
            .map(str::trim_end)
            .filter(|line| !line.is_empty())
            .collect();
        if lines.is_empty() {
            return Err(AocError::InvalidInput("grid must not be empty"));
        }
        let width = lines[0].chars().count();
        let mut tiles = Vec::with_capacity(lines.len());
        for line in &lines {
            if line.chars().count() != width {
                return Err(AocError::InvalidInput("every row must have the same width"));
            }
            let row = line
                .chars()
                .map(Tile::from_char)
                .collect::<Result<Vec<_>, _>>()?;
            tiles.push(row);
        }
        Ok(Grid { tiles })
    }

    fn rows(&self) -> usize {
        self.tiles.len()
    }

    fn cols(&self) -> usize {
        self.tiles[0].len()
    }

    fn start(&self) -> Result<(usize, usize), AocError> {
        self.tiles
            .iter()
            .enumerate()
            .find_map(|(r, row)| {
                row.iter()
                    .position(|tile| *tile == Tile::Start)
                    .map(|c| (r, c))
            })
            .ok_or(AocError::InvalidInput("grid has no start tile"))
    }

    fn step(&self, pos: (usize, usize), dir: Direction) -> Option<(usize, usize)> {
        let (dr, dc) = dir.delta();
        let r = pos.0 as isize + dr;
        let c = pos.1 as isize + dc;
        if r < 0 || c < 0 || r as usize >= self.rows() || c as usize >= self.cols() {
            return None;
        }
        Some((r as usize, c as usize))
    }

    fn tile_at(&self, pos: (usize, usize)) -> Tile {
        self.tiles[pos.0][pos.1]
    }

    /// Finds the directions the start tile actually connects in, by checking which
    /// neighbors have a pipe that connects back toward it.
    fn start_directions(&self, start: (usize, usize)) -> Vec<Direction> {
        ALL_DIRECTIONS
            .into_iter()
            .filter(|&dir| {
                self.step(start, dir)
                    .map(|next| self.tile_at(next).connects(dir.opposite()))
                    .unwrap_or(false)
            })
            .collect()
    }

    /// Walks the main loop starting at `S`, returning every tile it passes through
    /// in order. The start tile's real pipe shape is inferred from its neighbors.
    pub(crate) fn find_loop(&self) -> Result<Vec<(usize, usize)>, AocError> {
        let start = self.start()?;
        let start_dirs = self.start_directions(start);
        let [mut dir, _] = start_dirs.as_slice() else {
            return Err(AocError::InvalidInput(
                "start tile must connect to exactly two pipes",
            ));
        };

        let mut pos = start;
        let mut loop_tiles = Vec::new();
        loop {
            loop_tiles.push(pos);
            pos = self
                .step(pos, dir)
                .ok_or(AocError::InvalidInput("loop runs off the edge of the grid"))?;
            if pos == start {
                return Ok(loop_tiles);
            }
            let Tile::Pipe(a, b) = self.tile_at(pos) else {
                return Err(AocError::InvalidInput(
                    "loop passes through a non-pipe tile",
                ));
            };
            let incoming = dir.opposite();
            dir = if a == incoming {
                b
            } else if b == incoming {
                a
            } else {
                return Err(AocError::InvalidInput(
                    "pipe does not connect back to the previous tile",
                ));
            };
        }
    }
}

/// Counts tiles enclosed by a simple loop using the shoelace formula for its
/// area and Pick's theorem (`area = interior + boundary / 2 - 1`) to recover
/// the interior point count.
pub(crate) fn enclosed_area(loop_tiles: &[(usize, usize)]) -> i64 {
    let n = loop_tiles.len();
    let area2: i64 = (0..n)
        .map(|i| {
            let (r1, c1) = loop_tiles[i];
            let (r2, c2) = loop_tiles[(i + 1) % n];
            (c1 as i64) * (r2 as i64) - (c2 as i64) * (r1 as i64)
        })
        .sum();
    let area = area2.abs() / 2;
    let boundary = n as i64;
    area - boundary / 2 + 1
}
