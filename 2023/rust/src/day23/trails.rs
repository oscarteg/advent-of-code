use std::collections::HashMap;

use crate::day23::custom_error::AocError;

/// How slope tiles behave while hiking.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Slopes {
    /// Stepping onto a slope forces the next step downhill (part 1).
    Icy,
    /// Slopes are ordinary path tiles (part 2).
    Dry,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum Direction {
    North,
    East,
    South,
    West,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum Tile {
    Path,
    Forest,
    Slope(Direction),
}

/// A `(row, column)` tile coordinate.
type Position = (usize, usize);

const DIRECTIONS: [Direction; 4] = [
    Direction::North,
    Direction::East,
    Direction::South,
    Direction::West,
];

impl Direction {
    /// The neighbour in this direction, unless it lies past the top or left edge.
    fn step(self, (row, column): Position) -> Option<Position> {
        match self {
            Direction::North => Some((row.checked_sub(1)?, column)),
            Direction::East => Some((row, column + 1)),
            Direction::South => Some((row + 1, column)),
            Direction::West => Some((row, column.checked_sub(1)?)),
        }
    }
}

impl TryFrom<char> for Tile {
    type Error = AocError;

    fn try_from(tile: char) -> Result<Self, Self::Error> {
        match tile {
            '.' => Ok(Tile::Path),
            '#' => Ok(Tile::Forest),
            '^' => Ok(Tile::Slope(Direction::North)),
            '>' => Ok(Tile::Slope(Direction::East)),
            'v' => Ok(Tile::Slope(Direction::South)),
            '<' => Ok(Tile::Slope(Direction::West)),
            unknown => Err(AocError::UnknownTile(unknown)),
        }
    }
}

/// A trail map: a rectangular tile grid entered through the single path tile of
/// the top row and left through the single path tile of the bottom row.
#[derive(Debug)]
pub struct Map {
    tiles: Vec<Vec<Tile>>,
    start: Position,
    end: Position,
}

impl Map {
    pub fn parse(input: &str) -> Result<Self, AocError> {
        let tiles = input
            .lines()
            .map(str::trim_end)
            .filter(|line| !line.is_empty())
            .map(|line| {
                line.chars()
                    .map(Tile::try_from)
                    .collect::<Result<Vec<_>, _>>()
            })
            .collect::<Result<Vec<_>, _>>()?;
        let (Some(top), Some(bottom)) = (tiles.first(), tiles.last()) else {
            return Err(AocError::InvalidInput("empty map"));
        };
        if tiles.iter().any(|row| row.len() != top.len()) {
            return Err(AocError::InvalidInput("rows differ in length"));
        }
        let start = (
            0,
            path_column(top).ok_or(AocError::InvalidInput("no path tile in the top row"))?,
        );
        let end = (
            tiles.len() - 1,
            path_column(bottom).ok_or(AocError::InvalidInput("no path tile in the bottom row"))?,
        );
        Ok(Map { tiles, start, end })
    }

    /// The number of steps of the longest hike from the entrance to the exit.
    pub fn longest_hike(&self, slopes: Slopes) -> Result<usize, AocError> {
        self.compress(slopes).longest_hike().ok_or(AocError::NoHike)
    }

    /// Tiles outside the grid count as forest, which keeps every walk inside it.
    fn tile(&self, (row, column): Position) -> Tile {
        self.tiles
            .get(row)
            .and_then(|row| row.get(column))
            .copied()
            .unwrap_or(Tile::Forest)
    }

    /// The tiles reachable from `from` in a single step.
    fn steps(&self, from: Position, slopes: Slopes) -> impl Iterator<Item = Position> + '_ {
        let downhill = match (slopes, self.tile(from)) {
            (Slopes::Icy, Tile::Slope(direction)) => Some(direction),
            _ => None,
        };
        DIRECTIONS
            .into_iter()
            .filter(move |direction| downhill.is_none_or(|forced| *direction == forced))
            .filter_map(move |direction| direction.step(from))
            .filter(move |&to| self.tile(to) != Tile::Forest)
    }

    /// The entrance, the exit and every tile where the trails fork. Slopes never
    /// change which tiles fork, so both parts share one set of nodes.
    fn nodes(&self) -> Vec<Position> {
        let mut nodes = vec![self.start];
        if self.end != self.start {
            nodes.push(self.end);
        }
        for (row, tiles) in self.tiles.iter().enumerate() {
            for (column, tile) in tiles.iter().enumerate() {
                let position = (row, column);
                if *tile != Tile::Forest
                    && position != self.start
                    && position != self.end
                    && self.steps(position, Slopes::Dry).count() > 2
                {
                    nodes.push(position);
                }
            }
        }
        nodes
    }

    /// Replaces the corridors between the nodes by weighted edges.
    fn compress(&self, slopes: Slopes) -> Graph {
        let nodes = self.nodes();
        let indices: HashMap<Position, usize> = nodes
            .iter()
            .enumerate()
            .map(|(index, &position)| (position, index))
            .collect();
        let mut edges = vec![Vec::new(); nodes.len()];
        for (from, &position) in nodes.iter().enumerate() {
            for first in self.steps(position, slopes) {
                if let Some(edge) = self.follow(position, first, slopes, &indices) {
                    edges[from].push(edge);
                }
            }
        }
        // `nodes` lists the entrance first and the exit second; a single-row map
        // is entered and left through the same tile and so has one node.
        Graph {
            edges,
            start: 0,
            end: usize::from(self.end != self.start),
        }
    }

    /// Follows the corridor leaving a node through `first` up to the next node.
    /// A corridor that dead-ends, or that leads into a slope pointing back,
    /// yields no edge.
    fn follow(
        &self,
        node: Position,
        first: Position,
        slopes: Slopes,
        indices: &HashMap<Position, usize>,
    ) -> Option<Edge> {
        let (mut previous, mut current, mut steps) = (node, first, 1);
        loop {
            if let Some(&to) = indices.get(&current) {
                return Some(Edge { to, steps });
            }
            let next = self.steps(current, slopes).find(|&next| next != previous)?;
            (previous, current, steps) = (current, next, steps + 1);
        }
    }
}

fn path_column(row: &[Tile]) -> Option<usize> {
    row.iter().position(|tile| *tile == Tile::Path)
}

#[derive(Debug, Clone, Copy)]
struct Edge {
    to: usize,
    steps: usize,
}

/// The trail map with its corridors collapsed: the nodes are the entrance, the
/// exit and the forks, and every edge carries the steps along one corridor.
#[derive(Debug)]
struct Graph {
    edges: Vec<Vec<Edge>>,
    start: usize,
    end: usize,
}

impl Graph {
    fn longest_hike(&self) -> Option<usize> {
        let mut visited = vec![false; self.edges.len()];
        self.search(self.start, &mut visited)
    }

    /// Finding the longest path is NP-hard, so every hike that visits each node
    /// at most once is enumerated.
    fn search(&self, node: usize, visited: &mut [bool]) -> Option<usize> {
        if node == self.end {
            return Some(0);
        }
        visited[node] = true;
        let mut longest = None;
        for &Edge { to, steps } in &self.edges[node] {
            if !visited[to] {
                longest = longest.max(self.search(to, visited).map(|rest| steps + rest));
            }
        }
        visited[node] = false;
        longest
    }
}
