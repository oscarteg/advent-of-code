use anyhow::{ensure, Context};

type P1 = usize;
type P2 = usize;

#[derive(Clone)]
struct Grid {
    trees: Vec<Vec<u32>>,
    rows: usize,
    cols: usize,
}

impl Grid {
    fn parse(input: &str) -> anyhow::Result<Self> {
        let mut trees: Vec<Vec<u32>> = Vec::new();
        for line in input.lines() {
            let line = line.trim();
            if line.is_empty() {
                continue;
            }
            let mut row = Vec::new();
            for ch in line.chars() {
                let height = ch
                    .to_digit(10)
                    .with_context(|| format!("Invalid digit: {}", ch))?;
                row.push(height);
            }
            ensure!(!row.is_empty(), "Empty row in grid");
            if !trees.is_empty() {
                ensure!(
                    row.len() == trees[0].len(),
                    "Inconsistent row length: {} vs {}",
                    row.len(),
                    trees[0].len()
                );
            }
            trees.push(row);
        }
        ensure!(!trees.is_empty(), "Empty grid");
        let rows = trees.len();
        let cols = trees[0].len();
        Ok(Grid { trees, rows, cols })
    }

    fn is_visible(&self, r: usize, c: usize) -> bool {
        let height = self.trees[r][c];

        if (0..c).all(|col| self.trees[r][col] < height) {
            return true;
        }
        if ((c + 1)..self.cols).all(|col| self.trees[r][col] < height) {
            return true;
        }
        if (0..r).all(|row| self.trees[row][c] < height) {
            return true;
        }
        if ((r + 1)..self.rows).all(|row| self.trees[row][c] < height) {
            return true;
        }

        false
    }

    fn scenic_score(&self, r: usize, c: usize) -> usize {
        let height = self.trees[r][c];

        let left = (0..c)
            .rev()
            .position(|col| self.trees[r][col] >= height)
            .map(|p| p + 1)
            .unwrap_or(c);

        let right = ((c + 1)..self.cols)
            .position(|col| self.trees[r][col] >= height)
            .map(|p| p + 1)
            .unwrap_or(self.cols - c - 1);

        let up = (0..r)
            .rev()
            .position(|row| self.trees[row][c] >= height)
            .map(|p| p + 1)
            .unwrap_or(r);

        let down = ((r + 1)..self.rows)
            .position(|row| self.trees[row][c] >= height)
            .map(|p| p + 1)
            .unwrap_or(self.rows - r - 1);

        left * right * up * down
    }
}

pub fn solve(input: &str) -> anyhow::Result<(P1, P2)> {
    let grid = Grid::parse(input)?;

    let mut visible_count = 0;
    for r in 0..grid.rows {
        for c in 0..grid.cols {
            if grid.is_visible(r, c) {
                visible_count += 1;
            }
        }
    }

    let mut max_score = 0;
    for r in 0..grid.rows {
        for c in 0..grid.cols {
            let score = grid.scenic_score(r, c);
            max_score = max_score.max(score);
        }
    }

    Ok((visible_count, max_score))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_part1_example() {
        let input = include_str!("../../input/day8-example.txt");
        let (part1, _) = solve(input).unwrap();
        assert_eq!(part1, 21);
    }

    #[test]
    fn test_part2_example() {
        let input = include_str!("../../input/day8-example.txt");
        let (_, part2) = solve(input).unwrap();
        assert_eq!(part2, 8);
    }

    #[test]
    fn test_single_tree() {
        let input = "5\n";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 1);
        assert_eq!(part2, 0);
    }

    #[test]
    fn test_malformed_input() {
        let input = "12a34\n";
        assert!(solve(input).is_err());
    }

    #[test]
    fn test_malformed_inconsistent_cols() {
        let input = "123\n12\n";
        assert!(solve(input).is_err());
    }
}
