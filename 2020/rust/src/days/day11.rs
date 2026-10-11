use anyhow::Result;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum Cell {
    Floor,
    Empty,
    Occupied,
}

impl Cell {
    fn from_char(c: char) -> Result<Self> {
        match c {
            '.' => Ok(Cell::Floor),
            'L' => Ok(Cell::Empty),
            '#' => Ok(Cell::Occupied),
            _ => anyhow::bail!("invalid cell: {}", c),
        }
    }
}

#[derive(Debug, Clone)]
struct Grid {
    cells: Vec<Vec<Cell>>,
}

impl Grid {
    fn parse(input: &str) -> Result<Self> {
        let cells: Result<Vec<Vec<Cell>>> = input
            .lines()
            .filter(|line| !line.is_empty())
            .map(|line| line.chars().map(Cell::from_char).collect())
            .collect();
        Ok(Grid { cells: cells? })
    }

    fn rows(&self) -> usize {
        self.cells.len()
    }

    fn cols(&self) -> usize {
        if self.rows() == 0 {
            0
        } else {
            self.cells[0].len()
        }
    }

    fn get(&self, r: isize, c: isize) -> Option<Cell> {
        if r < 0 || c < 0 {
            return None;
        }
        let r = r as usize;
        let c = c as usize;
        if r < self.rows() && c < self.cols() {
            Some(self.cells[r][c])
        } else {
            None
        }
    }

    fn count_adjacent_occupied(&self, r: usize, c: usize) -> usize {
        let mut count = 0;
        for dr in -1..=1 {
            for dc in -1..=1 {
                if dr == 0 && dc == 0 {
                    continue;
                }
                if let Some(Cell::Occupied) = self.get(r as isize + dr, c as isize + dc) {
                    count += 1;
                }
            }
        }
        count
    }

    fn count_visible_occupied(&self, r: usize, c: usize) -> usize {
        let mut count = 0;
        for dr in -1..=1 {
            for dc in -1..=1 {
                if dr == 0 && dc == 0 {
                    continue;
                }
                let mut nr = r as isize + dr;
                let mut nc = c as isize + dc;
                loop {
                    match self.get(nr, nc) {
                        None => break,
                        Some(Cell::Floor) => {
                            nr += dr;
                            nc += dc;
                        }
                        Some(Cell::Empty) => break,
                        Some(Cell::Occupied) => {
                            count += 1;
                            break;
                        }
                    }
                }
            }
        }
        count
    }

    fn step(&self, threshold: usize, count_fn: fn(&Grid, usize, usize) -> usize) -> Grid {
        let mut new_cells = self.cells.clone();
        for (r, row) in new_cells.iter_mut().enumerate() {
            for (c, cell) in row.iter_mut().enumerate() {
                *cell = match self.cells[r][c] {
                    Cell::Floor => Cell::Floor,
                    Cell::Empty => {
                        if count_fn(self, r, c) == 0 {
                            Cell::Occupied
                        } else {
                            Cell::Empty
                        }
                    }
                    Cell::Occupied => {
                        if count_fn(self, r, c) >= threshold {
                            Cell::Empty
                        } else {
                            Cell::Occupied
                        }
                    }
                };
            }
        }
        Grid { cells: new_cells }
    }

    fn count_occupied(&self) -> usize {
        self.cells
            .iter()
            .flat_map(|row| row.iter())
            .filter(|&&cell| cell == Cell::Occupied)
            .count()
    }

    fn stabilize(&mut self, threshold: usize, count_fn: fn(&Grid, usize, usize) -> usize) {
        loop {
            let next = self.step(threshold, count_fn);
            if next.cells == self.cells {
                break;
            }
            self.cells = next.cells;
        }
    }
}

pub fn solve(input: &str) -> Result<(usize, usize)> {
    let mut grid1 = Grid::parse(input)?;
    grid1.stabilize(4, Grid::count_adjacent_occupied);
    let part1 = grid1.count_occupied();

    let mut grid2 = Grid::parse(input)?;
    grid2.stabilize(5, Grid::count_visible_occupied);
    let part2 = grid2.count_occupied();

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_example() {
        let input = include_str!("../../input/2020/day11-example.txt");
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 37);
        assert_eq!(part2, 26);
    }

    #[test]
    fn test_single_seat() {
        let input = "L";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 1);
        assert_eq!(part2, 1);
    }

    #[test]
    fn test_floor_only() {
        let input = ".";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 0);
        assert_eq!(part2, 0);
    }

    #[test]
    fn test_invalid_cell() {
        let input = "X";
        assert!(solve(input).is_err());
    }
}
