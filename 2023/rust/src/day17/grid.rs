use crate::day17::custom_error::AocError;

/// A rectangular grid of single-digit heat-loss costs.
pub struct Grid {
    costs: Vec<u32>,
    width: usize,
    height: usize,
}

#[derive(Clone, Copy, PartialEq, Eq, Hash, Debug)]
pub enum Direction {
    Up,
    Down,
    Left,
    Right,
}

impl Direction {
    pub const ALL: [Direction; 4] = [
        Direction::Up,
        Direction::Down,
        Direction::Left,
        Direction::Right,
    ];

    pub fn opposite(self) -> Direction {
        match self {
            Direction::Up => Direction::Down,
            Direction::Down => Direction::Up,
            Direction::Left => Direction::Right,
            Direction::Right => Direction::Left,
        }
    }

    fn delta(self) -> (i32, i32) {
        match self {
            Direction::Up => (0, -1),
            Direction::Down => (0, 1),
            Direction::Left => (-1, 0),
            Direction::Right => (1, 0),
        }
    }
}

impl Grid {
    pub fn parse(input: &str) -> Result<Self, AocError> {
        let rows: Vec<&str> = input.trim().lines().collect();
        if rows.is_empty() {
            return Err(AocError::InvalidInput("grid is empty"));
        }
        let width = rows[0].len();
        if width == 0 {
            return Err(AocError::InvalidInput("grid rows are empty"));
        }
        let mut costs = Vec::with_capacity(width * rows.len());
        for row in &rows {
            if row.len() != width {
                return Err(AocError::InvalidInput("grid rows have mismatched widths"));
            }
            for ch in row.chars() {
                let digit = ch
                    .to_digit(10)
                    .ok_or(AocError::InvalidInput("grid contains a non-digit"))?;
                costs.push(digit);
            }
        }
        Ok(Grid {
            costs,
            width,
            height: rows.len(),
        })
    }

    pub fn width(&self) -> usize {
        self.width
    }

    pub fn height(&self) -> usize {
        self.height
    }

    pub fn cost(&self, x: usize, y: usize) -> u32 {
        self.costs[y * self.width + x]
    }

    /// Steps one block from `(x, y)` in `direction`, or `None` if that leaves the grid.
    pub fn step(&self, x: usize, y: usize, direction: Direction) -> Option<(usize, usize)> {
        let (dx, dy) = direction.delta();
        let nx = x as i32 + dx;
        let ny = y as i32 + dy;
        if nx < 0 || ny < 0 || nx as usize >= self.width || ny as usize >= self.height {
            return None;
        }
        Some((nx as usize, ny as usize))
    }
}
